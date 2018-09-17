package payDataStruct

const (
	WxPayStatusWaitForWxPreId = "WaitForWxPreId"
	WxPayStatusWaitForUserPay = "WaitForUserPay"
	WxPayStatusSuccess        = "PaySuccess"
	WxPayStatusFinished       = "Finished"
)

// 微信公众号支付数据库结构
type WxPayBill struct {
	UserId        int    `json:"userId,string"`
	AppId         string `json:"appId"`
	MchId         string `json:"mchId"`
	BillId        string `json:"billId"`        // 商户订单号
	TransactionId string `json:"transactionId"` //微信支付订单号
	ProductId     string `json:"productId"`     //商户产品id
	OpenId        string `json:"openId"`
	TotalFee      int    `json:"totalFee,string"`   //订单价格
	Status        string `json:"status"`            // 订单状态
	CreateTime    int64  `json:"createTime,string"` // 订单创建时间
	FinishTime    int64  `json:"finishTime,string"` // 订单完成时间

}
