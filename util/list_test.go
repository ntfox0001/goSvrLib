package util_test

import (
	"fmt"
	"goSvrLib/util"
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	list := util.NewList()
	list.PushBack(reflect.SelectCase{})
	var s1 interface{}
	s1 = list.ToSlice()

	s2 := s1.([]reflect.SelectCase)

	fmt.Println(s2)
}
