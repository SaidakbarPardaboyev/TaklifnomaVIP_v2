package models

import "time"

type GetAllOrdersFilter struct {
	AccountID    string
	Status       *string
	TemplateCode *string
	FromDate     *time.Time
	ToDate       *time.Time
	Page         *int
	Limit        *int
	SortBy       *string
	Order        *string
}

type GetTransactionByPaymentIDRequest struct {
	PaymentID string
}

type GetTransactionByOrderIDRequest struct {
	OrderID string
}

type GetListTransactionRequest struct {
	AccountID  *string
	OrderID    *string
	PaymentIDs *[]string
	States     *[]int
	FromDate   *time.Time
	ToDate     *time.Time
	Page       *int
	Limit      *int
}
