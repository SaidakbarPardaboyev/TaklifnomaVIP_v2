package templete_mongo

import (
	"time"

	"saidakbar.origin/db/mongo/entity"
)

// region Create
type (
	CreateTempleteModel struct {
		OrganizationID string
		Name           string
		// TODO: add domain-specific create fields
	}

	CreateTempleteResult struct {
		Templete *entity.TempleteEntity
	}
)

// region Update
type (
	UpdateTempleteModel struct {
		ID             string
		OrganizationID string
		Name           string
		// TODO: add domain-specific update fields
	}

	UpdateTempleteResult struct {
		TempleteNotFound bool
		Templete         *entity.TempleteEntity
	}
)

// region Delete
type (
	DeleteTempleteModel struct {
		ID             string
		OrganizationID string
	}

	DeleteTempleteResult struct {
		TempleteNotFound bool
		Templete         *entity.TempleteEntity
	}
)

// region Get All
type (
	GetAllTempletesModel struct {
		OrganizationID string
		Skip           *int
		Limit          *int
		Search         *string
		StartDate      *time.Time
		EndDate        *time.Time
	}

	GetAllTempletesResult struct {
		Templetes []*entity.TempleteEntity
	}
)

// region Get By ID
type (
	GetTempleteByIDModel struct {
		ID             string
		OrganizationID string
	}

	GetTempleteByIDResult struct {
		TempleteNotFound bool
		Templete         *entity.TempleteEntity
	}
)

// region Get Count
type (
	GetTempleteCountModel struct {
		OrganizationID string
		Search         *string
		StartDate      *time.Time
		EndDate        *time.Time
	}

	GetTempleteCountResult struct {
		Count int64
	}
)
