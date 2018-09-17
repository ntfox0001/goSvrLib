package slHttpClient

import (
	"goSvrLib/network"
	"goSvrLib/selectCase/selectCaseInterface"
	"goSvrLib/util"

	"goSvrLib/log"
)

var _self *HttpClientManager

type HttpClientManager struct {
	goPool *util.GoroutinePool
}

type HttpClientResult struct {
	Body string
	Err  error
}

var (
	GoPoolSize int
	ExecSize   int
)

func Instance() *HttpClientManager {
	if _self == nil {
		_self = &HttpClientManager{}
		_self.goPool = util.NewGoPool("HttpClientManager", GoPoolSize, ExecSize)
	}
	return _self
}

func (*HttpClientManager) Initial() {

}
func (*HttpClientManager) Release() {
	_self.goPool.Release()

	log.Debug("HttpClientManager release")
}

func (*HttpClientManager) HttpGet(cb *selectCaseInterface.CallbackHandler, url string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpGet(url)
		rt := HttpClientResult{
			Body: rtStr,
			Err:  err,
		}
		if cb != nil {
			cb.SendReturnMsgNoReturn(rt)
		}

	}
	_self.goPool.Go(hp, nil)
}

func (*HttpClientManager) HttpPost(cb *selectCaseInterface.CallbackHandler, url string, content string, contentType string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpPost(url, content, contentType)
		rt := HttpClientResult{
			Body: rtStr,
			Err:  err,
		}
		if cb != nil {
			cb.SendReturnMsgNoReturn(rt)
		}

	}
	_self.goPool.Go(hp, nil)
}

func (*HttpClientManager) HttpPostByHeader(cb *selectCaseInterface.CallbackHandler, url string, content string, contentType string, header map[string]string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpPostByHeader(url, content, contentType, header)
		rt := HttpClientResult{
			Body: rtStr,
			Err:  err,
		}
		if cb != nil {
			cb.SendReturnMsgNoReturn(rt)
		}

	}
	_self.goPool.Go(hp, nil)
}
