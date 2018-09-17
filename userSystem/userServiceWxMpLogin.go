package userSystem

import (
	"goSvrLib/userSystem/userDefine"

	"goSvrLib/selectCase/selectCaseInterface"
	"io"
	"io/ioutil"
	"net/http"

	"goSvrLib/log"
	jsoniter "github.com/json-iterator/go"
)

// 独立的go协程
func (u *UserService) wxmpLoginProcess(w http.ResponseWriter, r *http.Request) {

	log.Debug("+ user http arrived.")
	s, _ := ioutil.ReadAll(r.Body)

	req := userDefine.WxMpLoginReq{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(s, &req); err != nil {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, err.Error())
		return
	}

	waitTokenChan := make(chan userDefine.GenerateTokenResp, 1)

	gtReq := userDefine.GenerateTokenReq{
		WaitTokenChan: waitTokenChan,
		WxMpLoginReq:  req,
	}

	u.userMgr.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("GenerateTokenReq", nil, gtReq))

	tokenResp := <-waitTokenChan

	if tokenResp.Token == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	loginResp := userDefine.WxMpLoginResp{
		MsgId:   "WxMpLoginResp",
		Token:   tokenResp.Token,
		UserId:  tokenResp.UserData.UserId,
		ErrorId: "0",
	}

	if s, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(loginResp); err != nil {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, err.Error())
	} else {

		io.WriteString(w, string(s))
	}

	log.Debug("- user http left.")
}
