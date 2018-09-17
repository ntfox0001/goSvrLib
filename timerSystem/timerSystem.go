package timerSystem

import (
	"container/list"
	"goSvrLib/selectCase"
	"goSvrLib/selectCase/selectCaseInterface"
	"reflect"
	"sync/atomic"
	"time"

	"goSvrLib/log"
)

type TimerSystem struct {
	selectLoop    *selectCase.SelectLoop
	timerItemList *list.List
	count         uint64
}

var _self *TimerSystem

const (
	TimerSystemTimerUpNotify = "TimerSystemTimerUpNotify"
)

// 实现一个1分钟级别的时间事件回调
func Instance() *TimerSystem {
	if _self == nil {
		_self = &TimerSystem{
			selectLoop:    selectCase.NewSelectLoop("TimerSystem", 10, 10),
			timerItemList: list.New(),
			count:         1,
		}
	}
	return _self
}

func (*TimerSystem) Initial() error {

	// 计算最近一个下一分钟0秒
	// t := time.Now().Add(time.Second * 60)
	// tStr := fmt.Sprintf("%d-%02d-%02d %02d:%02d:00", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
	// startTime, err := time.Parse("2006-01-02 15:04:05", tStr)
	// if err != nil {
	// 	return err
	// }
	now := time.Now()
	startTime := now.Unix() - int64(now.Second()) + 60

	delayTime := startTime - time.Now().Unix()
	log.Info("TimerSystem will runing", "delay", delayTime)

	_self.selectLoop.GetHelper().RegisterEvent("AddTimer", _self.addTimer)
	_self.selectLoop.GetHelper().RegisterEvent("DelTimer", _self.delTimer)
	go func() {
		time.Sleep(time.Second * time.Duration(delayTime))

		ticker := time.NewTicker(time.Second * 60)

		_self.selectLoop.GetHelper().AddSelectCase(reflect.ValueOf(ticker.C), _self.tickerCallback)

		// 触发一次
		_self.tickerCallback(nil)
	}()

	go _self.selectLoop.Run()

	return nil
}

func (*TimerSystem) Release() {
	_self.Release()
}
func (*TimerSystem) addTimer(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	item := msg.Content.(*TimerItem)
	_self.timerItemList.PushBack(item)
	return true
}
func (*TimerSystem) delTimer(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	id := msg.Content.(uint64)
	for i := _self.timerItemList.Front(); i != nil; i = i.Next() {
		if i.Value.(*TimerItem).id == id {
			_self.timerItemList.Remove(i)
			break
		}
	}
	return true
}

func (*TimerSystem) tickerCallback(data interface{}) bool {
	now := time.Now().Unix()
	t := time.Now()
	// 遍历所有注册的item，找到时间小于当前时间的，发送消息之后，删除
	for i := _self.timerItemList.Front(); i != nil; {
		ti := i.Value.(*TimerItem)
		if ti.time <= now {
			// 有啥调啥
			if ti.cb != nil {

				ti.cb.SendReturnMsgNoReturn(&t)
			}

			// 有啥调啥
			if ti.f != nil {
				t := time.Now()
				ti.f(ti.ud, &t)
			}

			// 是否是loop
			if ti.loop == false {
				e := i
				i = i.Next()
				_self.timerItemList.Remove(e)
			} else {
				ti.timeUp()
				i = i.Next()
			}
		} else {
			i = i.Next()
		}
	}
	return true
}

// 添加一个timer，thread safe
func (*TimerSystem) AddTimer(item *TimerItem) uint64 {
	id := _self.getNextId()
	item.id = id

	_self.selectLoop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("AddTimer", nil, item))

	return id
}

// 删除一个timer，thread safe
func (*TimerSystem) DelTimer(id uint64) {
	_self.selectLoop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("DelTimer", nil, id))
}
func (*TimerSystem) getNextId() uint64 {
	// 返回一个唯一值
	return atomic.AddUint64(&_self.count, 1)
}
