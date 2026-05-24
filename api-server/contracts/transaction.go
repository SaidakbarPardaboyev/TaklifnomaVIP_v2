package contracts

import (
	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/plugins"
)

type TransactionContract struct {
	ID          string  `json:"id"`
	PaymentID   string  `json:"payment_id"`
	Amount      float64 `json:"amount"`
	CreatedTime int64   `json:"created_time"`
	PerformTime int64   `json:"perform_time"`
	CancelTime  int64   `json:"cancel_time"`
	State       int     `json:"state"`
	Reason      *int    `json:"reason"`
	AccountID   string  `json:"account_id"`
	OrderID     string  `json:"order_id"`
}

func BuildTransactionContract(t *mysql_entity.TransactionModel) TransactionContract {
	return TransactionContract{
		ID:          t.ID,
		PaymentID:   t.PaymentID,
		Amount:      t.Amount,
		CreatedTime: plugins.ConverterToPaymeTimeFormat(t.CreatedTime),
		PerformTime: plugins.ConverterToPaymeTimeFormat(t.PerformTime),
		CancelTime:  plugins.ConverterToPaymeTimeFormat(t.CancelTime),
		State:       t.State,
		Reason:      t.Reason,
		AccountID:   t.AccountID,
		OrderID:     t.OrderID,
	}
}

func BuildTransactionListContract(txs []*mysql_entity.TransactionModel) []TransactionContract {
	result := make([]TransactionContract, 0, len(txs))
	for _, t := range txs {
		result = append(result, BuildTransactionContract(t))
	}
	return result
}
