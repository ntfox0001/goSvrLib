package slHttpClient

import (
	"goSvrLib/network"
	"goSvrLib/selectCase/selectCaseInterface"
)

// caCertPath:证书文件，用来验证服务器的证书是否真实
func (*HttpClientManager) HttpsGet(cb *selectCaseInterface.CallbackHandler, url, caCertPath string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpsGet(url, caCertPath)
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

// caCertPath:证书文件，用来验证服务器的证书是否真实
func (*HttpClientManager) HttpsPost(cb *selectCaseInterface.CallbackHandler, url, content, contentType, caCertPath string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpsPost(url, content, contentType, caCertPath)
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

// caCertPath:证书文件，用来验证服务器的证书是否真实
func (*HttpClientManager) HttpsPostByHeader(cb *selectCaseInterface.CallbackHandler, url, content, contentType string, header map[string]string, caCertPath string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpsPostByHeader(url, content, contentType, header, caCertPath)
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

// https two-way ---------------------------------------------------------------------------------------------------------

// caCertPath:证书文件，用来验证服务器的证书是否真实
// crtPath, keyPath: 用于验证客户端真实性的文件
func (*HttpClientManager) HttpsTwoWayGet(cb *selectCaseInterface.CallbackHandler, url, caCertPath, crtPath, keyPath string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpsTwoWayGet(url, caCertPath, crtPath, keyPath)
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

// caCertPath:证书文件，用来验证服务器的证书是否真实
// crtPath, keyPath: 用于验证客户端真实性的文件
func (*HttpClientManager) HttpsTwoWayPost(cb *selectCaseInterface.CallbackHandler, url, content, contentType, caCertPath, crtPath, keyPath string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpsTwoWayPost(url, content, contentType, caCertPath, crtPath, keyPath)
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

// caCertPath:证书文件，用来验证服务器的证书是否真实
// crtPath, keyPath: 用于验证客户端真实性的文件
func (*HttpClientManager) HttpsTwoWayPostByHeader(cb *selectCaseInterface.CallbackHandler, url, content, contentType string, header map[string]string, caCertPath, crtPath, keyPath string) {
	hp := func(data interface{}) {
		rtStr, err := network.SyncHttpsTwoWayPostByHeader(url, content, contentType, header, caCertPath, crtPath, keyPath)
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
