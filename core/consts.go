package core

import "time"

const (
	OtpTTL   = 2 * time.Minute
	TokenTTL = 30 * 24 * time.Hour
)

const (
	TransactionStateCreated               = 1
	TransactionStateCompleted             = 2
	TransactionStateCanceledWhenCreated   = -1
	TransactionStateCanceledWhenPerformed = -2

	TransactionCancelReasonWhenCreated   = 3
	TransactionCancelReasonWhenPerformed = 5
)
