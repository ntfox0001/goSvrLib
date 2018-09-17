package util_test

import (
	"fmt"
	"goSvrLib/util"
	"testing"
)

func TestUniqueId(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := util.GetUniqueId()
		fmt.Println(s, "  ", len(s))
	}

}
