package requestmodels

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type GeneratePaymeLinkRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

func (r GeneratePaymeLinkRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OrderID, validation.Required, is.UUID),
		validation.Field(&r.Amount, validation.Required),
	)
}

type PaymeWebhookRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

func (r PaymeWebhookRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Method, validation.Required),
	)
}

type PaymeAccount struct {
	AccountID string `json:"account_id"`
	OrderID   string `json:"order_id"`
}

type CheckPerformTransactionRequest struct {
	Account PaymeAccount `json:"account"`
	Amount  float64      `json:"amount"`
}

func (r CheckPerformTransactionRequest) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Amount, validation.Required),
	); err != nil {
		return err
	}
	return validation.ValidateStruct(&r.Account,
		validation.Field(&r.Account.OrderID, validation.Required, is.UUID),
		validation.Field(&r.Account.AccountID, validation.Required, is.UUID),
	)
}

type CreateTransactionRequest struct {
	Account   PaymeAccount `json:"account"`
	PaymentID string       `json:"id"`
	Amount    int          `json:"amount"`
	TimePayme int64        `json:"time"`
}

func (r CreateTransactionRequest) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.PaymentID, validation.Required),
		validation.Field(&r.Amount, validation.Required),
		validation.Field(&r.TimePayme, validation.Required),
	); err != nil {
		return err
	}
	return validation.ValidateStruct(&r.Account,
		validation.Field(&r.Account.OrderID, validation.Required, is.UUID),
		validation.Field(&r.Account.AccountID, validation.Required, is.UUID),
	)
}

type PerformTransactionRequest struct {
	ID string `json:"id"`
}

func (r PerformTransactionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
	)
}

type CancelTransactionRequest struct {
	ID     string `json:"id"`
	Reason int    `json:"reason"`
}

func (r CancelTransactionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
	)
}

type CheckTransactionRequest struct {
	ID string `json:"id"`
}

func (r CheckTransactionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
	)
}

type GetStatementRequest struct {
	To   int64 `json:"to"`
	From int64 `json:"from"`
}

func (r GetStatementRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.To, validation.Required),
		validation.Field(&r.From, validation.Required),
	)
}
