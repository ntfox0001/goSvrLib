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

type UserServiceParams struct {
	Listenip          string
	Port              string
	CertFile, KeyFile string
	Servcb            usInterface.IServiceCallback
	UsrMgrcb          usInterface.IUserCallback
	Usrcb             usInterface.IUserCallback
}

func NewUserService(params UserServiceParams) *UserService {

	usrServ := UserService{
		userMgr:  NewUserManager(params.Listenip, params.Port, params.UsrMgrcb, params.Usrcb),
		listenip: params.Listenip,
		port:     params.Port,
		server:   nil,
		ssl:      true,
		certFile: params.CertFile,
		keyFile:  params.KeyFile,
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
