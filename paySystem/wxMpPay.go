package paySystem

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"goSvrLib/database"
	"goSvrLib/network"
	"goSvrLib/paySystem/payDataStruct"
	"goSvrLib/selectCase/selectCaseInterface"
	"goSvrLib/util"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"

	"goSvrLib/log"
)

const (
	WxPreBillUrl     = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	WxNotifyFailResp = `<xml>
	<return_code><![CDATA[FAIL]]></return_code>
	<return_msg><![CDATA[INVALIDFORMAT]]></return_msg>
  </xml>`
	WxNotifySuccessResp = `<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
</xml>`
)

type beginPayReq struct {
	req payDataStruct.WxUnifiedorderReq
	cb  *selectCaseInterface.CallbackHandler
}

type WxMpPay struct {
	appId     string
	mchId     string
	mckKey    string
	notifyUrl string
	count     int64
	goPool    *util.GoroutinePool
}

// 微信公众号支付
// 微信公众号处理流程
// 客户端通知服务器发起支付-》服务器调用“统一下单”获得支付preid-》客户端收到preid，呼叫微信支付sdk
// -》用户输入密码确认-》微信后台收到确认-》通知服务器结果，服务器处理结果-》发给客户端显示结果
// 这里解决的问题是服务器和微信服务器和数据库订单状态管理
// 监听端口用来告诉微信服务器发送数据到哪里
func NewWxMpPay(appId string, mchId string, mckKey string, notifyUrl string) *WxMpPay {
	wp := &WxMpPay{
		appId:     appId,
		mchId:     mchId,
		mckKey:    mckKey,
		count:     10232,
		notifyUrl: notifyUrl,
	}

	return wp
}
func (w *WxMpPay) generateBill(pd payDataStruct.ClientWxPayBase) string {
	// 产生 一个全局唯一字符串
	c := atomic.AddInt64(&w.count, 1)
	s := fmt.Sprintf("%s%s%d%s%s%s%d", w.appId, w.mchId, c, w.mckKey, pd.OpenId, pd.ProductId, time.Now().Unix())
	h := md5.New()
	h.Write([]byte(s))
	md5Str := hex.EncodeToString(h.Sum(nil))

	return md5Str
}

//发起一笔用户微信支付
func (w *WxMpPay) BeginPay(pd payDataStruct.ClientWxPayReq) (*payDataStruct.ClientWxPayResp, error) {
	req := payDataStruct.WxUnifiedorderReq{
		ClientWxPayBase: pd.ClientWxPayBase,
		AppId:           w.appId,
		MchId:           w.mchId,
		Key:             w.mckKey,
		NonceStr:        "wxpay",
		OutTradeNo:      w.generateBill(pd.ClientWxPayBase),
		TradeType:       "JSAPI",
		NotifyUrl:       w.notifyUrl,
	}

	// 在数据库创建订单,发送到数据库保存
	op := database.Instance().NewOperation("call WxPayBillInsertPreBill(?,?,?,?,?,?,?,?,?)", pd.UserId, req.AppId, req.MchId, req.OutTradeNo, 0, req.ProductId, req.OpenId, req.TotalFee, payDataStruct.WxPayStatusWaitForWxPreId)
	_, err := database.Instance().SyncExecOperation(op)
	if err != nil {
		log.Error("sql:WxPayBill_Insert error", "err", err.Error())
		return nil, err
	}

	// 创建签名数据
	if xmlStr, err := util.MakeWxSign(req); err != nil {
		log.Error("MakeWxSign error", "err", err.Error())
		return nil, err
	} else {
		//向微信服务器发送“统一下单”请求
		wxRespStr, err := network.SyncHttpPost(WxPreBillUrl, xmlStr, network.ContentTypeText)
		fmt.Println(wxRespStr)
		if err != nil {
			log.Error("SyncHttpPost error", "err", err.Error())
			return nil, err
		}
		// 解析微信返回值
		var resp payDataStruct.WxUnifiedorderResp
		if err := xml.Unmarshal([]byte(wxRespStr), &resp); err != nil {
			// 如果格式解析失败，那是一个严重错误
			log.Error("Failed to Unmarshal of WxMpPay resp.", "resp", wxRespStr)
			return nil, err
		}

		// 检查微信返回值, 通信成功标识
		if resp.ReturnCode != "SUCCESS" {
			// 返回微信错误，写数据库关闭订单
			op := database.Instance().NewOperation("call WxPayBill_UpdateStatusByBillId(?,?)", req.OutTradeNo, resp.ReturnCode)
			if _, err := database.Instance().SyncExecOperation(op); err != nil {
				log.Error("sql:WxPayBill_UpdateStatusByBillId error", "resp", wxRespStr)
				return nil, err
			}
			log.Warn("wx pay failed", "resp", resp)
			return nil, err
		}

		// 检查微信返回值，订单成功标识
		if resp.ResultCode != "SUCCESS" {
			// 返回微信错误，写数据库关闭订单
			op := database.Instance().NewOperation("call WxPayBill_UpdateStatusByBillId(?,?)", req.OutTradeNo, resp.ResultCode)
			if _, err := database.Instance().SyncExecOperation(op); err != nil {
				log.Error("sql:WxPayBill_UpdateStatusByBillId error", "resp", wxRespStr)
				return nil, err
			}
			log.Warn("wx pay failed.", "resp", resp)
			return nil, err
		}

		// 都成功了
		// 更新数据库
		op := database.Instance().NewOperation("call WxPayBill_UpdateStatusByBillId(?,?)", req.OutTradeNo, payDataStruct.WxPayStatusWaitForUserPay)
		if _, err := database.Instance().SyncExecOperation(op); err != nil {
			log.Error("sql:WxPayBill_UpdateStatusByBillId error", "resp", wxRespStr)
			return nil, err
		}
		// 返回resp,等待用户付款，快输入密码，快快快~
		clientResp := payDataStruct.ClientWxPayResp{
			WxUnifiedorderResp: resp,
			UserId:             pd.UserId,
			BillId:             req.OutTradeNo,
		}
		return &clientResp, err

	}

}

// 微信支付通知
func (*WxMpPay) wxNotifyReq(w http.ResponseWriter, r *http.Request) {
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

		notify := payDataStruct.PaySystemNotify{
			ExtentData: req,
			ProductId:  wxpaybill.ProductId,
			UserId:     wxpaybill.UserId,
			PayType:    payDataStruct.PaySystemNotify_PayType_Wx,
		}
		log.Info("wxPay success", "data", notify)
		// 发送回调
		_self.wxCallback.SendReturnMsgNoReturn(notify)

		// 返回成功
		io.WriteString(w, WxNotifySuccessResp)
	}

}
