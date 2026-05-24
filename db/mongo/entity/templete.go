package entity

// TempleteEntity is the example MongoDB document.
// Rename this file and struct to match your domain (e.g. ContractorEntity, OrderEntity).
type TempleteEntity struct {
	BaseEntity     `bson:",inline"`
	OrganizationID string `json:"organization_id" bson:"organization_id"`
	Name           string `json:"name" bson:"name"`
	// TODO: add domain-specific fields here
}
