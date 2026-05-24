package entity

// Field name constants — used when building MongoDB filter/update bson.M documents.
// Add one block per collection your service owns.
const (
	FieldTempleteID             = "_id"
	FieldTempleteOrganizationID = "organization_id"
	FieldTempleteName           = "name"
	FieldTempleteIsDeleted      = "is_deleted"
	FieldTempleteDeletedAt      = "deleted_at"
	FieldTempleteCreatedAt      = "created_at"
	FieldTempleteUpdatedAt      = "updated_at"
)
