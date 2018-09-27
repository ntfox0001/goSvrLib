package payDataStruct

import "goSvrLib/database/dbtools/dbtoolsData"

const (
	WxPayStatusWaitForWxPreId = "WaitForWxPreId"
	WxPayStatusWaitForUserPay = "WaitForUserPay"
	WxPayStatusSuccess        = "PaySuccess"
	WxPayStatusFinished       = "Finished"
)

// 微信支付数据库结构
type WxPayBill struct {
	WxPayBillTable dbtoolsData.TableName
	UserId         int    `json:"userId,string" dbdef:"int"`
	AppId          string `json:"appId" dbdef:"varchar(128)"`
	MchId          string `json:"mchId" dbdef:"varchar(128)"`
	BillId         string `json:"billId" dbdef:"varchar(256),prim"`   // 商户订单号
	TransactionId  string `json:"transactionId" dbdef:"varchar(256)"` //微信支付订单号
	ProductId      string `json:"productId" dbdef:"varchar(32)"`      //商户产品id
	OpenId         string `json:"openId" dbdef:"varchar(128)"`
	TotalFee       int    `json:"totalFee,string" dbdef:"int"`   //订单价格
	Status         string `json:"status" dbdef:"tinyint"`        // 订单状态
	CreateTime     int64  `json:"createTime,string" dbdef:"int"` // 订单创建时间
	FinishTime     int64  `json:"finishTime,string" dbdef:"int"` // 订单完成时间

}
