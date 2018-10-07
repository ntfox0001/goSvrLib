package paySystem

import (
	"goSvrLib/log"
	"goSvrLib/network"

	jsoniter "github.com/json-iterator/go"
)

const (
	sandUrl = "https://sandbox.itunes.apple.com/verifyReceipt"
	buyUrl  = "https://buy.itunes.apple.com/verifyReceipt"
)

type iapPayDataReq struct {
	Receipt_data             string `json:"receipt-data"`
	Password                 string `json:"password"`
	Exclude_old_transactions string `json:"exclude-old-transactions"`
}

type iapPayDataResp struct {
	Status                      string `json:"status"`
	Receipt                     string `json:"receipt"`
	Latest_receipt              string `json:"latest_receipt"`
	Latest_receipt_info         string `json:"latest_receipt_info"`
	Latest_expired_receipt_info string `json:"latest_expired_receipt_info"`
	Pending_renewal_info        string `json:"pending_renewal_info"`
	Is_retryable                string `json:"is-retryable"`
}

func (ps *PaySystem) AppleIAPPay(receipt string, productId string) {
	req := iapPayDataReq{
		Receipt_data: receipt,
	}

	reqStr, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(req)
	if err != nil {
		log.Warn("req marshal failed.", "req", reqStr)
		return
	}

	sandRespStr, err := network.SyncHttpPost(sandUrl, reqStr, network.ContentTypeJson)
	if err != nil {

	}

}
