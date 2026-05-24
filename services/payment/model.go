package payment_service

import (
	"saidakbar.origin/db/mongo/entity"
	mysql_entity "saidakbar.origin/db/mysql/entity"
)

type GeneratePaymeLinkModel struct {
	AccountID string
	OrderID   string
	Amount    float64
}

type GeneratePaymeLinkResult struct {
	Link                 *string
	OrderNotFound        bool
	Unauthorized         bool
	IncorrectAmount      bool
}

type CheckPerformTransactionModel struct {
	AccountID string
	OrderID   string
	Amount    float64
}

type CheckPerformTransactionResult struct {
	Succeed              bool
	Order                *entity.OrderEntity
	OrderNotFound        bool
	AccountNotFound      bool
	Unauthorized         bool
	IncorrectAmount      bool
	AlreadyPaid          bool
	ErrorCode            int
	MessageUz, MessageRu, MessageEn string
}

type CreateTransactionModel struct {
	AccountID string
	OrderID   string
	PaymentID string
	Amount    int
	TimePayme int64
}

type CreateTransactionResult struct {
	Succeed              bool
	Transaction          *mysql_entity.TransactionModel
	OrderNotFound        bool
	AccountNotFound      bool
	Unauthorized         bool
	IncorrectAmount      bool
	AlreadyPaid          bool
	TransactionExists    bool
	ErrorCode            int
	MessageUz, MessageRu, MessageEn string
}

type PerformTransactionModel struct {
	PaymentID string
}

type PerformTransactionResult struct {
	Succeed             bool
	Transaction         *mysql_entity.TransactionModel
	TransactionNotFound bool
	OrderNotFound       bool
	ErrorCode           int
	MessageUz, MessageRu, MessageEn string
}

type CancelTransactionModel struct {
	PaymentID string
	Reason    int
}

type CancelTransactionResult struct {
	Succeed             bool
	Transaction         *mysql_entity.TransactionModel
	TransactionNotFound bool
	OrderNotFound       bool
	ErrorCode           int
	MessageUz, MessageRu, MessageEn string
}

type CheckTransactionModel struct {
	PaymentID string
}

type CheckTransactionResult struct {
	Succeed             bool
	Transaction         *mysql_entity.TransactionModel
	TransactionNotFound bool
	OrderNotFound       bool
	ErrorCode           int
	MessageUz, MessageRu, MessageEn string
}

type GetStatementModel struct {
	From int64
	To   int64
}

type GetStatementResult struct {
	Transactions []*mysql_entity.TransactionModel
}
