package paySystem

import (
	"fmt"
	"goSvrLib/commonError"
	"goSvrLib/log"
	"goSvrLib/network"
	"goSvrLib/paySystem/payDataStruct"
	"goSvrLib/util"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	sandUrl  = "https://sandbox.itunes.apple.com/verifyReceipt"
	buyUrl   = "https://buy.itunes.apple.com/verifyReceipt"
	netError = 1
)

type applePayItem struct {
	UserId        int
	Receipt       string
	ProductId     string
	RepeatCount   int // 重试次数
	updateManager *util.UpdateManager
}

func newApplePayItem(userId int, receipt string, productId string) *applePayItem {
	item := &applePayItem{
		UserId:      userId,
		Receipt:     receipt,
		ProductId:   productId,
		RepeatCount: 0,
	}

	item.updateManager = util.NewUpdateManager2("applePay", time.Second, item.update)

	return item
}

func (i *applePayItem) update() {
	if i.RepeatCount
}

func (i *applePayItem) beginPay() {
	// 保存订单数据
	billId := util.GetUniqueId()
	if err := _self.PayRecord_NewBill(i.UserId, billId, i.ProductId, 0, "APPLE", i.Receipt, false); err != nil {
		log.Warn("apple NewBill failed.", "receipt", i.Receipt, "userId", i.UserId)
		return
	}

	// 获得结果
	if resp, err := i.getReceiptResp(i.Receipt); err != nil {
		log.Warn(err.Error(), "userId", i.UserId)
		if err.(commonError.CommError).GetType() == netError {

		}
		return
	} else {
		// 检查
		i.validatingReceipt(i.UserId, billId, i.ProductId, resp)
	}
}

// 验证收据是否有效
func (i *applePayItem) getReceiptResp(receipt string) (payDataStruct.IapPayDataResp, error) {
	req := payDataStruct.IapPayDataReq{
		Receipt_data: receipt,
	}

	reqStr, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(req)
	if err != nil {
		log.Warn("req marshal failed.", "req", reqStr)
		return payDataStruct.IapPayDataResp{}, commonError.NewCommErr("req marshal failed.", 0)
	}

	sandRespStr, err := network.SyncHttpPost(sandUrl, reqStr, network.ContentTypeJson)
	if err != nil {
		log.Warn("sandbox http post failed", "err", err.Error())
		return payDataStruct.IapPayDataResp{}, commonError.NewCommErr("sandbox http post failed", netError)
	}

	resp := payDataStruct.IapPayDataResp{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(sandRespStr, &resp); err != nil {
		log.Warn("the format of sandbox's resp failed.")
		return payDataStruct.IapPayDataResp{}, commonError.NewCommErr("the format of sandbox's resp failed.", 0)
	}

	// check status
	if resp.Status == 21007 {
		// 这是一个正式服的收据
		buyRespStr, err := network.SyncHttpPost(buyUrl, reqStr, network.ContentTypeJson)
		if err != nil {
			log.Warn("http post failed", "err", err.Error())
			return payDataStruct.IapPayDataResp{}, commonError.NewCommErr("http post failed.", netError)
		}

		resp = payDataStruct.IapPayDataResp{}
		if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(buyRespStr, &resp); err != nil {
			log.Warn("the format of resp failed.")
			return payDataStruct.IapPayDataResp{}, commonError.NewCommErr("the format of resp failed.", 0)
		}
	}

	return resp, nil
}

// 验证收据是否正确
func (i *applePayItem) validatingReceipt(userId int, billId string, productId string, resp payDataStruct.IapPayDataResp) {
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
