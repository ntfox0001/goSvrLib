package paySystem

import (
	"goSvrLib/commonError"
	"goSvrLib/database"
	"goSvrLib/paySystem/payDataStruct"
	"goSvrLib/util"
	"time"
)

/*
	支付数据库操作
	插入预支付id记录
	插入支付记录
	更新支付记录状态(4个状态)

	这里的函数都是同步函数
*/

// 创建新订单，等待用户支付，一般用于微信支付等
// 需要预先向第三方支付服务器申请预支付id，
// 用于客户端拉起第三方支付程序
// extentInfo 可以保存appid或者receipt
func (*PaySystem) PayRecord_NewBill(userId int, billId string, productId string, fee int, billType string, extentInfo string) error {
	// inuserId int,inbillId varchar(256),intransactionId text,inproductId varchar(32),intotalFee int,instatus varchar(32),increateTime int,infinishTime int,inextentInfo text,inbillType varchar(32)
	op := database.Instance().NewOperation("call PayBillTable_Insert(?,?,?,?,?,?,?,?,?,?)",
		userId, billId, "", productId, fee, payDataStruct.PayStatusWaitForUserPay, time.Now().Unix(), 0, extentInfo, billType)

	_, err := database.Instance().SyncExecOperation(op)

	return err
}

// 设置订单支付成功，下一步回调逻辑层，完成订单
func (*PaySystem) PayRecord_SetPayStatusSuccess(billId string, transactionId string) error {
	//inbillId varchar(256), intransactionId varchar(256),instatus tinyint
	op := database.Instance().NewOperation("call PayBillTable_PaySuccess(?,?,?)", billId, transactionId, payDataStruct.PayStatusSuccess)

	_, err := database.Instance().SyncExecOperation(op)
	return err
}

// 设置一个订单状态为错误，错误信息不能超过32个字符
func (*PaySystem) PayRecord_SetError(billId string, errInfo string) error {
	op := database.Instance().NewOperation("call PayBillTable_UpdateStatusByBillId(?,?)", billId, errInfo)
	_, err := database.Instance().SyncExecOperation(op)
	return err
}

// 完成订单
func (*PaySystem) PayRecord_Finish(billId string) error {
	// inbillId varchar(256),instatus varchar(32)
	op := database.Instance().NewOperation("call PayBillTable_UpdateStatusByBillId(?,?)", billId, payDataStruct.PayStatusFinished)
	_, err := database.Instance().SyncExecOperation(op)
	return err
}

// 查询订单
func (*PaySystem) PayRecord_Query(billId string) (payDataStruct.PayBillData, error) {
	op := database.Instance().NewOperation("call PayBillTable_QueryByBillId(?)", billId)
	rt, err := database.Instance().SyncExecOperation(op)
	if err != nil {
		return payDataStruct.PayBillData{}, err
	}
	payBillDS := rt.FirstSet()
	if len(payBillDS) != 1 {
		return payDataStruct.PayBillData{}, commonError.NewStringErr2("bill dataset length must be 1. len:", len(payBillDS))
	}

	pd := payDataStruct.PayBillData{}
	if err := util.I2Stru(payBillDS[0], &pd); err != nil {
		return payDataStruct.PayBillData{}, commonError.NewStringErr2("bill data format error")
	}

	return pd, nil
}
