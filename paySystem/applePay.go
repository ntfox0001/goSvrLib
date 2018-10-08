package paySystem

import (
	"fmt"
	"goSvrLib/commonError"
	"goSvrLib/log"
	"goSvrLib/network"
	"goSvrLib/paySystem/payDataStruct"
	"goSvrLib/util"

	jsoniter "github.com/json-iterator/go"
)

const (
	sandUrl = "https://sandbox.itunes.apple.com/verifyReceipt"
	buyUrl  = "https://buy.itunes.apple.com/verifyReceipt"
)

func (ps *PaySystem) BeginPay(userId int, receipt string, productId string) {

	// 保存订单数据
	billId := util.GetUniqueId()
	if err := _self.PayRecord_NewBill(userId, billId, productId, 0, "APPLE", receipt, false); err != nil {
		log.Warn("apple NewBill failed.", "receipt", receipt, "userId", userId)
		return
	}

	if resp, err := getReceiptResp(receipt); err != nil {
		log.Warn(err.Error(), "userId", userId)
		return
	} else {
		validatingReceipt(userId, billId, productId, resp)
	}
}

// 验证收据是否有效
func getReceiptResp(receipt string) (payDataStruct.IapPayDataResp, error) {
	req := payDataStruct.IapPayDataReq{
		Receipt_data: receipt,
	}

	reqStr, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(req)
	if err != nil {
		log.Warn("req marshal failed.", "req", reqStr)
		return payDataStruct.IapPayDataResp{}, commonError.NewStringErr("req marshal failed.")
	}

	sandRespStr, err := network.SyncHttpPost(sandUrl, reqStr, network.ContentTypeJson)
	if err != nil {
		log.Warn("sandbox http post failed", "err", err.Error())
		return payDataStruct.IapPayDataResp{}, commonError.NewStringErr("sandbox http post failed")
	}

	resp := payDataStruct.IapPayDataResp{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(sandRespStr, &resp); err != nil {
		log.Warn("the format of sandbox's resp failed.")
		return payDataStruct.IapPayDataResp{}, commonError.NewStringErr("the format of sandbox's resp failed.")
	}

	// check status
	if resp.Status == 21007 {
		// 这是一个正式服的收据
		buyRespStr, err := network.SyncHttpPost(buyUrl, reqStr, network.ContentTypeJson)
		if err != nil {
			log.Warn("http post failed", "err", err.Error())
			return payDataStruct.IapPayDataResp{}, commonError.NewStringErr("http post failed.")
		}

		resp = payDataStruct.IapPayDataResp{}
		if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(buyRespStr, &resp); err != nil {
			log.Warn("the format of resp failed.")
			return payDataStruct.IapPayDataResp{}, commonError.NewStringErr("the format of resp failed.")
		}
	}

	return resp, nil
}

// 验证收据是否正确
func validatingReceipt(userId int, billId string, productId string, resp payDataStruct.IapPayDataResp) {
	if resp.Status == 0 {
		transaction_id := ""
		for _, v := range resp.Receipt.In_app {
			if productId == v.Product_id {
				transaction_id = v.Transaction_id
			}
		}
		if err := _self.PayRecord_SetPayStatusSuccess(billId, transaction_id); err != nil {
			log.Error("apple SetPayStatusSuccess failed", "err", err.Error())
		}
	} else {
		if err := _self.PayRecord_SetError(billId, fmt.Sprint(resp.Status)); err != nil {
			log.Error("apple SetError failed", "err", err.Error())
		}
	}
}
