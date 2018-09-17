package paySystem

import (
	"encoding/xml"
	"fmt"
	"goSvrLib/commonError"
	"goSvrLib/database"
	"goSvrLib/network"
	"goSvrLib/paySystem/payDataStruct"
	"goSvrLib/selectCase/selectCaseInterface"
	"goSvrLib/util"
	"io"
	"io/ioutil"
	"net/http"
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
	_self.server.RegisterRouter(path, network.RouterHandler{ProcessHttpFunc: _self.wxNotifyReq})

	pay := NewWxMpPay(appId, mchId, mckKey, wxNotifyUrl)

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

// 微信支付通知
func (*PaySystem) wxNotifyReq(w http.ResponseWriter, r *http.Request) {
	s, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(s))
	req := payDataStruct.WxPayNotifyReq{}
	if err := xml.Unmarshal(s, &req); err != nil {
		log.Error("wxNotifyReq", "Marshal error", err.Error())
		io.WriteString(w, WxNotifyFailResp)
		return
	}

	// 查找数据库，是否有这个订单
	op := database.Instance().NewOperation("call WxPayBill_Query(?)", req.OutTradeNo)
	if rt, err := database.Instance().SyncExecOperation(op); err != nil {
		log.Error("WxNotifyReq sql:WxPayBill_Query", "err", err.Error())
		io.WriteString(w, WxNotifyFailResp)
		return
	} else {
		wxPayBillDS := rt.FirstSet()
		if len(wxPayBillDS) != 1 {
			log.Error("WxPayBill does not found", "bill", req.OutTradeNo, "req", req)
			io.WriteString(w, WxNotifyFailResp)
			return
		}
		// 解析数据
		var wxpaybill payDataStruct.WxPayBill
		for _, v := range wxPayBillDS {
			if err := util.I2Stru(v, &wxpaybill); err != nil {
				log.Error("Invalid to WxPayBill format")
				io.WriteString(w, WxNotifyFailResp)
				return
			}
			break
		}

		// 订单状态必须是等待用户支付
		if wxpaybill.Status != payDataStruct.WxPayStatusWaitForUserPay {
			if wxpaybill.Status == payDataStruct.WxPayStatusFinished {
				// 多余的补单，直接忽略
				return
			}
			log.Error("Invalid to WxPayBill status", "billId", wxpaybill.BillId, "status", wxpaybill.Status)
			io.WriteString(w, WxNotifyFailResp)
			return
		}

		// 逻辑层处理完成，那么关闭这个订单
		op := database.Instance().NewOperation("call WxPayBill_Finish(?,?,?)", wxpaybill.BillId, req.TransactionId, payDataStruct.WxPayStatusFinished)
		if _, err := database.Instance().SyncExecOperation(op); err != nil {
			log.Error("sql:WxPayBill_UpdateStatus", "err", err.Error())
			io.WriteString(w, WxNotifyFailResp)
			return
		}

		// 调用通知函数
		if _self.wxCallback == nil {
			log.Error("WxCallback is nil")
			io.WriteString(w, WxNotifyFailResp)
			return
		}

		notify := payDataStruct.PaySystemWxNotify{
			WxPayNotifyReq: req,
			ProductId:      wxpaybill.ProductId,
			UserId:         wxpaybill.UserId,
		}
		log.Info("wxPay success", "data", notify)
		// 发送回调
		_self.wxCallback.SendReturnMsgNoReturn(notify)

		// 返回成功
		io.WriteString(w, WxNotifySuccessResp)
	}

}
