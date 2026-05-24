package transaction_service

import mysql_entity "saidakbar.origin/db/mysql/entity"

type GetByOrderIDModel struct {
	OrderID   string
	AccountID string
}

type GetByOrderIDResult struct {
	Transaction  *mysql_entity.TransactionModel
	NotFound     bool
	Unauthorized bool
}

type GetListModel struct {
	AccountID string
	Page      *int
	Limit     *int
}

type GetListResult struct {
	Transactions []*mysql_entity.TransactionModel
}
