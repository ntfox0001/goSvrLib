package database

import (
	"container/list"
	"fmt"
	"goSvrLib/commonError"

	"goSvrLib/log"
)

type PrepareData struct {
	dataSet   *list.List
	dataCount int
}

func NewPrepareData() *PrepareData {
	return &PrepareData{
		dataSet:   list.New(),
		dataCount: -1,
	}
}

func (p *PrepareData) AddData(args ...interface{}) error {
	if p.dataCount != -1 {
		if len(args) != p.dataCount {
			errStr := fmt.Sprintf("args count error have: %d, want: %d", len(args), p.dataCount)
			log.Debug(errStr)
			return commonError.NewStringErr(errStr)
		}
	} else {
		p.dataCount = len(args)
	}
	p.dataSet.PushBack(args)
	return nil
}
