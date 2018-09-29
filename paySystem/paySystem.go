package paySystem

import (
	"goSvrLib/commonError"
	"goSvrLib/network"
	"goSvrLib/paySystem/payDataStruct"
	"goSvrLib/selectCase/selectCaseInterface"
	"goSvrLib/util"
	"net/url"

	"goSvrLib/log"
)

type PaySystem struct {
	wxNotifyUrl string
	server      *network.Server
	wxMpPayMap  map[string]*WxMpPay
	goPool      *util.GoroutinePool
	wxCallback  *selectCaseInterface.CallbackHandler
}

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

// 设置一个微信支付回调，回调的参数是一个WxPayNotifyReq
func (*PaySystem) SetWxCallbackFunc(wxCallback *selectCaseInterface.CallbackHandler) {
	_self.wxCallback = wxCallback

}

// 添加wx支付数据，wxCallback当微信服务器返回支付成功时，PaySystem会先验证消息，更新数据库，然后调用这个函数
func (*PaySystem) AddWxPay(appId string, mchId string, mckKey string, wxNotifyUrl string) error {
	if _, ok := _self.wxMpPayMap[appId]; ok {
		return commonError.NewStringErr("There is AppId already:" + appId)
	}
	u, err := url.Parse(wxNotifyUrl)
	if err != nil {
		return err
	}
	path := u.EscapedPath()
	if path == "" {
		return commonError.NewStringErr("EscapedPath is empty:" + wxNotifyUrl)
	}

	pay := NewWxMpPay(appId, mchId, mckKey, wxNotifyUrl)
	_self.server.RegisterRouter(path, network.RouterHandler{ProcessHttpFunc: pay.wxNotifyReq})

	_self.wxMpPayMap[appId] = pay

	return nil
}

// 发起一笔微信支付,通过cb，返回一个*ClientWxPayResp类型数据
func (*PaySystem) WxPay(pd payDataStruct.ClientWxPayReq, cb *selectCaseInterface.CallbackHandler) error {

	if pay, ok := _self.wxMpPayMap[pd.AppId]; !ok {
		return commonError.NewStringErr("appid does not exist:" + pd.AppId)
	} else {
		_self.goPool.Go(func(data interface{}) {
			resp, err := pay.BeginPay(pd)
			if err != nil {
				cb.SendReturnMsgNoReturn(&payDataStruct.ClientWxPayResp{ErrorId: err.Error()})
			} else {
				cb.SendReturnMsgNoReturn(resp)
			}
		}, nil)

	}
	return nil
}
