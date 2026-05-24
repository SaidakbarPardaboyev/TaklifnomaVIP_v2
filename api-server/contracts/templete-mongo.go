package contracts

import (
	"saidakbar.origin/db/mongo/entity"
	"saidakbar.origin/plugins"
)

type TempleteMongoContract struct {
	ID             string  `json:"id"`
	OrganizationID string  `json:"organization_id"`
	Name           string  `json:"name"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	IsDeleted      bool    `json:"is_deleted"`
	DeletedAt      *string `json:"deleted_at"`
}

func CreateTempleteMongoContract(t *entity.TempleteEntity) TempleteMongoContract {
	if t == nil {
		return TempleteMongoContract{}
	}

	return TempleteMongoContract{
		ID:             t.ID.Hex(),
		OrganizationID: t.OrganizationID,
		Name:           t.Name,
		CreatedAt:      plugins.FormatTime(t.CreatedAt),
		UpdatedAt:      plugins.FormatTime(t.UpdatedAt),
		IsDeleted:      t.IsDeleted,
		DeletedAt:      plugins.FormatTimePtr(t.DeletedAt),
	}
}
