package contracts

import (
	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/plugins"
)

type TemplateMysqlContract struct {
	ID             string  `json:"id"`
	OrganizationID string  `json:"organization_id"`
	Name           string  `json:"name"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	IsDeleted      bool    `json:"is_deleted"`
	DeletedAt      *string `json:"deleted_at"`
}

func CreateTemplateMysqlContract(t *mysql_entity.TempleteModel) TemplateMysqlContract {
	if t == nil {
		return TemplateMysqlContract{}
	}

	return TemplateMysqlContract{
		ID:             t.ID,
		OrganizationID: t.OrganizationID,
		Name:           t.Name,
		CreatedAt:      plugins.FormatTime(t.CreatedAt),
		UpdatedAt:      plugins.FormatTime(t.UpdatedAt),
		IsDeleted:      t.IsDeleted,
		DeletedAt:      plugins.FormatTimePtr(t.DeletedAt),
	}
}
