package mysql_entity

// Column name constants used in GORM queries.
// Add one block per table your service owns.
const (
	FieldAccountID        = "id"
	FieldAccountFullName  = "full_name"
	FieldAccountPhone     = "phone"
	FieldAccountChatID    = "chat_id"
	FieldAccountIsActive  = "is_active"
	FieldAccountCreatedAt = "created_at"
	FieldAccountUpdatedAt = "updated_at"
	FieldAccountDeletedAt = "deleted_at"
	FieldAccountIsDeleted = "is_deleted"

	FieldTransactionID          = "id"
	FieldTransactionPaymentID   = "payment_id"
	FieldTransactionAmount      = "amount"
	FieldTransactionCreatedTime = "created_time"
	FieldTransactionPerformTime = "perform_time"
	FieldTransactionCancelTime  = "cancel_time"
	FieldTransactionState       = "state"
	FieldTransactionReason      = "reason"
	FieldTransactionTimePayme   = "time_payme"
	FieldTransactionAccountID   = "account_id"
	FieldTransactionOrderID     = "order_id"
)
