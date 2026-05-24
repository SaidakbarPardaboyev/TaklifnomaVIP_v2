package mysql_entity

// Column name constants used in GORM queries.
// Add one block per table your service owns.
const (
	FieldTempleteID             = "id"
	FieldTempleteOrganizationID = "organization_id"
	FieldTempleteName           = "name"
	FieldTempleteCreatedAt      = "created_at"
	FieldTempleteUpdatedAt      = "updated_at"
	FieldTempleteIsDeleted      = "is_deleted"
	FieldTempleteDeletedAt      = "deleted_at"

	FieldAccountID        = "id"
	FieldAccountFullName  = "full_name"
	FieldAccountPhone     = "phone"
	FieldAccountChatID    = "chat_id"
	FieldAccountIsActive  = "is_active"
	FieldAccountCreatedAt = "created_at"
	FieldAccountUpdatedAt = "updated_at"
	FieldAccountDeletedAt = "deleted_at"
	FieldAccountIsDeleted = "is_deleted"
)
