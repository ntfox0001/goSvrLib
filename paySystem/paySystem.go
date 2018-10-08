package paySystem

import (
	"fmt"
	"goSvrLib/commonError"
	"goSvrLib/network"
	"goSvrLib/paySystem/payDataStruct"
	"goSvrLib/selectCase/selectCaseInterface"
	"goSvrLib/util"

	"goSvrLib/log"
)

type PaySystem struct {
	wxNotifyUrl string
	server      *network.Server
	wxMpPayMap  map[string]*WxMpPay
	goPool      *util.GoroutinePool
	wxCallback  *selectCaseInterface.CallbackHandler
}

const (
	WxNotifyPath = "wxPayNotify"
)

var _self *PaySystem

func Instance() *PaySystem {
	if _self == nil {
		_self = &PaySystem{
			wxMpPayMap: make(map[string]*WxMpPay),
		}
	}
	return _self
}

func (*PaySystem) Initial(listenip, port string, goPoolSize, execSize int) error {

	_self.server = network.NewServer(listenip, port)
	_self.goPool = util.NewGoPool("PaySystem", goPoolSize, execSize)

	// 接受微信支付服务器通知的地址
	_self.wxNotifyUrl = fmt.Sprintf("http://%s:%s/%s", listenip, port, WxNotifyPath)
	return nil
}
func (*PaySystem) Release() {
	_self.goPool.Release()
}
func (*PaySystem) Run() {

	if _self.wxCallback == nil {
		log.Warn("PaySystem wxCallback is nil.")
	}
	_self.server.Start()
}

// 设置一个微信支付回调，回调的参数是一个 PaySystemNotify，当支付成功时调用
func (*PaySystem) SetWxCallbackFunc(wxCallback *selectCaseInterface.CallbackHandler) {
	_self.wxCallback = wxCallback

}

// 添加wx支付数据，wxCallback当微信服务器返回支付成功时，PaySystem会先验证消息，更新数据库，然后调用这个函数
func (*PaySystem) AddWxPay(appId string, mchId string, mckKey string) error {
	if _, ok := _self.wxMpPayMap[appId]; ok {
		return commonError.NewStringErr("There is AppId already:" + appId)
	}

	pay := NewWxMpPay(appId, mchId, mckKey, _self.wxNotifyUrl)
	_self.server.RegisterRouter(WxNotifyPath, network.RouterHandler{ProcessHttpFunc: wxNotifyReq})
	_self.wxMpPayMap[appId] = pay

	return nil
}

// 发起一笔微信支付,通过cb，返回一个prePayId string
func (*PaySystem) WxPay(pd payDataStruct.WxPayReqData, cb *selectCaseInterface.CallbackHandler) error {

	if pay, ok := _self.wxMpPayMap[pd.AppId]; !ok {
		return commonError.NewStringErr("appid does not exist:" + pd.AppId)
	} else {
		_self.goPool.Go(func(data interface{}) {
			resp, err := pay.BeginPay(pd)
			if err != nil {
				log.Warn("wx pay failed.", "err", err.Error())
			} else {
				cb.SendReturnMsgNoReturn(resp)
			}
		}, nil)

	}
	return nil
}

func (*PaySystem) ApplePay(userId int, receipt string, cb *selectCaseInterface.CallbackHandler) {
	_self.goPool.Go(func(data interface{}) {

		//_self.BeginPay()
	}, nil)
}
