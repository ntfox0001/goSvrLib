package util

import (
	"goSvrLib/log"
)

type UpdateManager struct {
	name       string
	count      uint64
	updateFunc map[uint64][]func()
}

func NewUpdateManager(name string) *UpdateManager {
	updateMgr := &UpdateManager{
		name:       name,
		count:      0,
		updateFunc: make(map[uint64][]func()),
	}

	return updateMgr
}

func (u *UpdateManager) Update() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("UpdateManager update error", "name", u.name, "err", err.(error).Error())
		}
	}()

	for k, v := range u.updateFunc {
		if k == 1 || u.count%k == 0 {
			for _, f := range v {
				f()
			}
		}
	}
	u.count++
}

// 添加一个更新项，这个更新不支持撤销
func (u *UpdateManager) Add(interval uint64, f func()) {
	a, ok := u.updateFunc[interval]
	if !ok {
		u.updateFunc[interval] = make([]func(), 0, 4)
	}
	u.updateFunc[interval] = append(a, f)
}
