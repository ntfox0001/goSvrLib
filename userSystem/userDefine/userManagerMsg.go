package userDefine

import (
	"crypto/rand"
	"crypto/sha256"
	"goSvrLib/network/networkInterface"
	"strconv"

	"github.com/btcsuite/btcutil/base58"
)

const (
	tokenPreName = "uguessuguess?!haha!"
)

type UserPair struct {
	Ac      networkInterface.IMsgHandler
	UnionId string
}

func NewToken(str string) string {
	b := make([]byte, 2)
	rand.Read(b)

	hashstr := tokenPreName + strconv.Itoa(int(b[1])) + str + strconv.Itoa(int(b[0]))

	h := sha256.New()
	h.Write([]byte(hashstr))
	rt := base58.CheckEncode(h.Sum(nil), 0)

	return rt
}

type FindTokenReq struct {
	Token      string
	WaitWxChan chan string
}

func NewFindTokenReq(token string, waitWxChan chan string) FindTokenReq {
	return FindTokenReq{
		Token:      token,
		WaitWxChan: waitWxChan,
	}
}

type WxMpLoginReq struct {
	UserData
}

type GenerateTokenReq struct {
	WaitTokenChan chan GenerateTokenResp
	UserData
}

type GenerateTokenResp struct {
	Token string
	UserData
}

type WxMpLoginResp struct {
	MsgId   string `json:"msgId"`
	Token   string `json:"token"`
	UserId  int    `json:"userId,string"`
	ErrorId string `json:"errorId"`
}

type NewUserInfoReq struct {
	UserData
	WaitTokenChan chan GenerateTokenResp
}

// wx code login
type WxMpCodeLoginReq struct {
	Code string `json:"code"`
}

type WxMpCodeLoginResp struct {
	
}