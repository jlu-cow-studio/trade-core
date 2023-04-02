package biz

import (
	"fmt"
	"testing"

	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/sanity-io/litter"
)

func TestRecharge(t *testing.T) {
	mysql.Init()

	transaction, err := Recharge("322", 200000)
	litter.Dump(transaction)
	litter.Dump(err)
}

func TestOrder(t *testing.T) {
	mysql.Init()

	transaction, order, err := Order("322", "1", 4)
	litter.Dump(transaction)
	litter.Dump(order)
	fmt.Println(err)
}

func TestOrderList(t *testing.T) {
	mysql.Init()

	list, err := OrderList("322", 0, 2)
	litter.Dump(list)
	fmt.Println(err)
}
