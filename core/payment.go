package core

import "fmt"

const (
	PaymeApiKeyPrefix = "Paycom:"

	// payme methods
	CheckPerformTransaction = "CheckPerformTransaction"
	CreateTransaction       = "CreateTransaction"
	PerformTransaction      = "PerformTransaction"
	CancelTransaction       = "CancelTransaction"
	CheckTransaction        = "CheckTransaction"
	GetStatement            = "GetStatement"

	// payme error codes
	PaymeAccessDenied           = -32504
	PaymeBadRequest             = -31099
	PaymeErrorNotFound          = -31050
	PaymeIncorrectAmount        = -31001
	PaymeErrorCannotBeCompleted = -31008
	PaymeInternalServerError    = -32400
	PaymeTransactionExists      = -31051

	PaymeReceiptTypeSale   = 0
	PaymeReceiptTypeReturn = 1
)

type Language int

const (
	uz Language = 1
	ru Language = 2
	en Language = 3
)

var (
	PaymeErrorMessageUserNotFound = map[Language]string{
		uz: "user not found: user_id = %s",
		ru: "user not found: user_id = %s",
		en: "user not found: user_id = %s",
	}
	PaymeErrorMessageOrderNotFound = map[Language]string{
		uz: "order not found: order_id = %s",
		ru: "order not found: order_id = %s",
		en: "order not found: order_id = %s",
	}
	PaymeErrorMessageOrderNotBelongToUser = map[Language]string{
		uz: "order not belong to user: user_id = %s, order_id = %s",
		ru: "order not belong to user: user_id = %s, order_id = %s",
		en: "order not belong to user: user_id = %s, order_id = %s",
	}
	PaymeErrorMessageIncorrectAmount = map[Language]string{
		uz: "incorrect amount: given = %f, expected = %f",
		ru: "incorrect amount: given = %f, expected = %f",
		en: "incorrect amount: given = %f, expected = %f",
	}
	PaymeErrorMessageOrderAlreadyPaid = map[Language]string{
		uz: "order was already paid",
		ru: "order was already paid",
		en: "order was already paid",
	}
	PaymeErrorMessageTransactionNotFound = map[Language]string{
		uz: "transaction not found: payment_id = %s",
		ru: "transaction not found: payment_id = %s",
		en: "transaction not found: payment_id = %s",
	}
	PaymeErrorMessageTransactionExists = map[Language]string{
		uz: "transaction already exists",
		ru: "transaction already exists",
		en: "transaction already exists",
	}
)

func CreatePaymeMessage(model map[Language]string, args ...any) (uzStr, ruStr, enStr string) {
	return fmt.Sprintf(model[uz], args...), fmt.Sprintf(model[ru], args...), fmt.Sprintf(model[en], args...)
}
