package mysql_entity

import "time"

type TransactionModel struct {
	ID          string     `gorm:"column:id;primaryKey;size:36"`
	RowOrder    int        `gorm:"column:row_order"`
	PaymentID   string     `gorm:"column:payment_id;size:255"`
	Amount      float64    `gorm:"column:amount;not null"`
	CreatedTime *time.Time `gorm:"column:created_time"`
	PerformTime *time.Time `gorm:"column:perform_time"`
	CancelTime  *time.Time `gorm:"column:cancel_time"`
	State       int        `gorm:"column:state"`
	Reason      *int       `gorm:"column:reason"`
	TimePayme   int64      `gorm:"column:time_payme"`
	AccountID   string     `gorm:"column:account_id;size:36;not null"`
	OrderID     string     `gorm:"column:order_id;size:255;not null"`
}

func (TransactionModel) TableName() string { return "transactions" }
