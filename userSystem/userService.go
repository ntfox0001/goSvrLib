package userSystem

import (
	"goSvrLib/network"
	"goSvrLib/userSystem/usInterface"

	"goSvrLib/log"
)

type UserService struct {
	userMgr           *UserManager
	listenip          string
	port              string
	server            *network.Server
	ssl               bool
	certFile, keyFile string
	callback          usInterface.IServiceCallback
}

func NewUserService(listenip string, port string, servcb usInterface.IServiceCallback, usrMgrcb usInterface.IUserCallback, usrcb usInterface.IUserCallback) *UserService {

	return NewUserServiceSsl(listenip, port, "", "", servcb, usrMgrcb, usrcb)
}

func NewUserServiceSsl(listenip string, port string, certFile, keyFile string, servcb usInterface.IServiceCallback, usrMgrcb usInterface.IUserCallback, usrcb usInterface.IUserCallback) *UserService {

	usrServ := UserService{
		userMgr:  NewUserManager(listenip, port, usrMgrcb, usrcb),
		listenip: listenip,
		port:     port,
		server:   nil,
		ssl:      true,
		certFile: certFile,
		keyFile:  keyFile,
	}

	return &usrServ
}

func (u *UserService) Initial() {
	if u.ssl {
		u.server = network.NewServerSsl(u.listenip, u.port, u.certFile, u.keyFile)
	} else {
		u.server = network.NewServer(u.listenip, u.port)
	}

	// 用户长连接
	wsr := network.NewRouterWSHandler(u.userMgr)
	wsr.DisableCheckOrigin(false)
	u.server.RegisterRouter("/user", wsr)

	// 注册微信消息
	wxmpRouter := network.RouterHandler{
		ProcessHttpFunc: u.wxmpLoginProcess,
	}
	u.server.RegisterRouter("/wxmpLogin", wxmpRouter)

	// 初始化回调
	if err := u.callback.Initial(u.server); err != nil {
		log.Error("callback initial error", "err", err.Error())
		return
	}

	if err := u.server.Start(); err != nil {
		log.Error("http server error", "err", err.Error())
		return
	}
}

func (u *UserService) Release() {
	u.callback.Release()
	u.server.Close()
	u.userMgr.Release()

	log.Debug("UserService release.")
}
