package userSystem

import (
	"goSvrLib/database"
	"goSvrLib/selectCase/selectCaseInterface"
	"goSvrLib/userSystem/userDefine"
)

type UserInfo struct {
	usrData *userDefine.UserData
}

// 异步加载 for userManagerManager
func asyncNewUserInfo(cb *selectCaseInterface.CallbackHandler, callbackMsgName string, req userDefine.WxMpLoginReq, userData interface{}) {

	op := database.Instance().NewOperation("call userInsert(?,?,?,?,?)", req.UnionId, req.OpenId, req.Nickname, req.Headimgurl, req.RefreshToken)
	op.UserData = userData
	database.Instance().ExecOperationForCB(cb, op)

	return
}

// 创建新用户
func newUserInfoForUserData(ud *userDefine.UserData) *UserInfo {
	ui := &UserInfo{
		usrData: ud,
	}

	return ui
}

func (u *UserInfo) GetUserData() *userDefine.UserData {
	return u.usrData
}
