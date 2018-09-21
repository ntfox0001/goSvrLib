package usInterface

import (
	"goSvrLib/network"
	"goSvrLib/selectCase/selectCaseInterface"
)

type IUserCallback interface {
	Initial(helper selectCaseInterface.ISelectLoopHelper) error
	Release()
}
type IServiceCallback interface {
	// 
	Initial(server *network.Server) error
	Release()
}
