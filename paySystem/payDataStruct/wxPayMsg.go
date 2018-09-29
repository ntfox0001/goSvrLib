package payDataStruct

import (
	"encoding/xml"
)

// 微信统一下单请求，发送时使用了json序列化，所以这里必须标注json tag
type WxUnifiedorderReq struct {
	ClientWxPayBase
	AppId      string `json:"appid"`        //公众账号ID
	MchId      string `json:"mch_id"`       //商户号
	NonceStr   string `json:"nonce_str"`    //随机字符串
	OutTradeNo string `json:"out_trade_no"` //商户订单号
	NotifyUrl  string `json:"notify_url"`   //通知地址
	TradeType  string `json:"trade_type"`   //交易类型
	Key        string `json:"key"`          // 商户key
}

// 请求成功后，会返回这个对象
type WxUnifiedorderResp struct {
	XmlName    xml.Name `json:"-" xml:"xml"`
	ReturnCode string   `json:"returnCode" xml:"return_code"`            //返回状态码,必然返回的
	RetrunMsg  string   `json:"returnMsg" xml:"return_msg"`              //返回信息,必然返回的
	AppId      string   `json:"appId" xml:"appid,omitempty"`             //公众账号ID
	MchId      string   `json:"mchId" xml:"mch_id,omitempty"`            //商户号
	DeviceInfo string   `json:"deviceInfo" xml:"device_info,omitempty"`  //设备号
	NonceStr   string   `json:"nonceStr" xml:"nonce_str,omitempty"`      //随机字符串
	ResultCode string   `json:"resultCode" xml:"result_code,omitempty"`  //业务结果
	ErrCode    string   `json:"errCode" xml:"err_code,omitempty"`        //错误代码
	ErrCodeDes string   `json:"errCodeDes" xml:"err_code_des,omitempty"` //错误代码描述
	TradeType  string   `json:"tradeType" xml:"trade_type,omitempty"`    //交易类型
	PrePayId   string   `json:"prePayId" xml:"prepay_id"`                //预支付交易会话标识
}

type ClientWxPayBase struct {
	DeviceInfo     string `json:"device_info"`      //设备号 可不填
	Attach         string `json:"attach"`           //附加数据 可不填
	Body           string `json:"body"`             //商品描述 可不填
	TotalFee       int    `json:"total_fee,string"` //标价金额 必填
	SpBillCreateIp string `json:"spbill_create_ip"` //终端IP 必填
	ProductId      string `json:"product_id"`       //商品ID 必填
	OpenId         string `json:"openid"`           //用户标识 必填
}

// 客户端发起发起微信支付请求所需数据
type ClientWxPayReq struct {
	ClientWxPayBase
	UserId int    `json:"userId,string"` // 玩家id 必填
	AppId  string `json:"appId"`         // 支付使用的appid
}

// 返回客户端的数据
type ClientWxPayResp struct {
	WxUnifiedorderResp
	UserId  int    `json:"userId,string"` // 玩家id
	BillId  string `json:"billId"`        // 商户订单号，用于查询订单是否支付成功
	ErrorId string `json:"errorId"`       //用来返回程序中其他错误
}

//----------------------------------------------------------------------------------------------------------------

// 微信服务器发送过来的数据
type WxPayNotifyReq struct {
	XmlName       xml.Name `xml:"xml"`
	ReturnCode    string   `xml:"return_code"` //返回状态码,必然返回的
	RetrunMsg     string   `xml:"return_msg"`  //返回信息,必然返回的
	ErrCode       string   `xml:"err_code"`
	AppId         string   `xml:"appid,omitempty"`       //公众账号ID
	MchId         string   `xml:"mch_id,omitempty"`      //商户号
	DeviceInfo    string   `xml:"device_info,omitempty"` //设备号
	NonceStr      string   `xml:"nonce_str,omitempty"`   //随机字符串
	OpenId        string   `xml:"openid,omitempty"`
	TotalFee      int      `xml:"total_fee"`
	TransactionId string   `xml:"transaction_id"` // 微信订单id
	OutTradeNo    string   `xml:"out_trade_no"`   // 商户订单号
	Attach        string   `xml:"attach"`         // 附加数据
}

type WxPayNotifyResp struct {
	XmlName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"` // SUCCESS
	ReturnMsg  string   `xml:"return_msg"`  // OK
}

// 支付系统回调函数的参数
type PaySystemNotify struct {
	ExtentData interface{}
	UserId     int    `json:"userId"`
	ProductId  string `json:"productId"` //商户产品id
	PayType    int    // 1：wx,2： apple
}

const (
	PaySystemNotify_PayType_Wx    = 1
	PaySystemNotify_PayType_Apple = 2
)
