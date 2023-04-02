package biz

import (
	"fmt"
	"testing"

	"github.com/jlu-cow-studio/common/dal/mysql"
)

func TestRecharge(t *testing.T) {
	mysql.Init()

	fmt.Println(Recharge("322", 10))
}
