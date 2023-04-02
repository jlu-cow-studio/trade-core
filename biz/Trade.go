package biz

import (
	"fmt"
	"log"
	"time"

	"github.com/jlu-cow-studio/common/dal/mysql"
	mysql_model "github.com/jlu-cow-studio/common/model/dao_struct/mysql"
	"github.com/sanity-io/litter"
	"gorm.io/gorm"
)

func Recharge(userId string, money float64) (*mysql_model.Transaction, error) {

	var transaction = new(mysql_model.Transaction)

	err := mysql.GetDBConn().Transaction(func(tx *gorm.DB) error {

		wallet := &mysql_model.Wallet{}

		if err := tx.Set("gorm:query_option", "FOR UPDATE").Table("wallet").Where("user_id = ?", userId).First(wallet).Error; err != nil {
			return err
		}

		if err := tx.Table("wallet").Where("user_id = ?", userId).UpdateColumn("balance", wallet.Balance+money).Error; err != nil {
			return err
		}

		transaction = &mysql_model.Transaction{
			WalletID:  wallet.ID,
			Type:      "deposit",
			Amount:    money,
			CreatedAt: time.Now(),
		}

		if err := tx.Table("wallet_transaction").Create(transaction).Error; err != nil {
			return err
		}

		return nil
	})
	return transaction, err
}

func Order(userId, itemId string, quantity int) (*mysql_model.Transaction, *mysql_model.Order, error) {

	var order = new(mysql_model.Order)
	var transaction = new(mysql_model.Transaction)

	err := mysql.GetDBConn().Transaction(func(tx *gorm.DB) error {
		//查询商品，扣减库存，扣减余额，生成订单，生成交易明细

		//获取商品信息
		item := &mysql_model.Item{}
		itemCount := new(int64)
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Table("items").Where("id = ?", itemId).Find(item).Count(itemCount).Error; err != nil {
			return err
		}
		if *itemCount != 1 {
			return fmt.Errorf("error when get item %v from db", itemId)
		}
		log.Printf("Get item %v from db\n", litter.Sdump(item))

		//获取钱包信息
		wallet := &mysql_model.Wallet{}
		walletCount := new(int64)
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Table("wallet").Where("user_id = ?", userId).Find(wallet).Count(walletCount).Error; err != nil {
			return err
		}
		if *walletCount != 1 {
			return fmt.Errorf("error when get wallet of %v from db", userId)
		}
		log.Printf("Get wallet %v from db\n", litter.Sdump(wallet))

		//计算库存
		if item.Stock < int32(quantity) {
			return fmt.Errorf("insufficient stock of item %v, %v needed but %v left", itemId, quantity, item.Stock)
		}
		log.Printf("The stock of item %v is sufficient\n", itemId)

		//计算价格
		if item.Price*float64(quantity) > wallet.Balance {
			return fmt.Errorf("insufficient balance of user %v", userId)
		}
		log.Printf("The balance of user %v is sufficient\n", userId)

		//扣减库存
		if err := tx.Table("items").Where("id = ?", itemId).UpdateColumn("stock", item.Stock-int32(quantity)).Error; err != nil {
			return err
		}
		log.Printf("Update the stock of item %v to %v\n", itemId, item.Stock-int32(quantity))

		//扣减余额
		if err := tx.Table("wallet").Where("user_id = ?", userId).UpdateColumn("balance", wallet.Balance-item.Price*float64(quantity)).Error; err != nil {
			return err
		}
		log.Printf("Update the balance of user %v to %v\n", userId, wallet.Balance-item.Price*float64(quantity))

		//生成订单
		order = &mysql_model.Order{
			UserID:           wallet.ID,
			ItemID:           int(item.ID),
			PurchaseQuantity: quantity,
			UnitPrice:        item.Price,
			TotalPrice:       item.Price * float64(quantity),
			PurchaseDate:     time.Now(),
		}
		if err := tx.Table("order_table").Create(order).Error; err != nil {
			return err
		}
		log.Printf("Create order %v for user %v and item %v\n", order.OrderID, userId, itemId)

		//生成交易记录
		transaction = &mysql_model.Transaction{
			WalletID: wallet.ID,
			Type:     "purchase",
			Amount:   item.Price * float64(quantity),
			OrderID:  order.OrderID,
		}

		if err := tx.Table("wallet_transaction").Create(transaction).Error; err != nil {
			return err
		}

		return nil
	})
	return transaction, order, err
}
