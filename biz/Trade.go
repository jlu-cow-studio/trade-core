package biz

import (
	"time"

	"github.com/jlu-cow-studio/common/dal/mysql"
	"gorm.io/gorm"
)

func Recharge(userId string, money float64) error {
	return mysql.GetDBConn().Transaction(func(tx *gorm.DB) error {

		wallet := &struct {
			ID        int       `gorm:"column:id"`
			UserID    int       `gorm:"column:user_id"`
			Balance   float64   `gorm:"column:balance"`
			CreatedAt time.Time `gorm:"column:created_at"`
			UpdatedAt time.Time `gorm:"column:updated_at"`
		}{}

		if err := tx.Set("gorm:query_option", "FOR UPDATE").Table("wallet").Where("user_id = ?", userId).First(wallet).Error; err != nil {
			return err
		}

		if err := tx.Table("wallet").Where("user_id = ?", userId).UpdateColumn("balance", wallet.Balance+money).Error; err != nil {
			return err
		}

		transaction := &struct {
			ID          uint64    `gorm:"primary_key;column:id"`
			WalletID    uint64    `gorm:"column:wallet_id"`
			Type        string    `gorm:"type:ENUM('deposit', 'purchase');column:type"`
			Amount      float64   `gorm:"type:decimal(10,2);column:amount"`
			OrderID     uint64    `gorm:"column:order_id"`
			Description string    `gorm:"type:varchar(255);column:description"`
			CreatedAt   time.Time `gorm:"column:created_at"`
		}{
			WalletID: uint64(wallet.ID),
			Type:     "deposit",
			Amount:   money,
		}

		if err := tx.Table("wallet_transaction").Create(transaction).Error; err != nil {
			return err
		}

		return nil
	})
}
