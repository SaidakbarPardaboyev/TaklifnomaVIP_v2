package templete_mysql

import (
	"time"

	mysql_entity "saidakbar.origin/db/mysql/entity"
)

// region Create
type (
	CreateTempleteModel struct {
		OrganizationID string
		Name           string
		// TODO: add domain-specific create fields
	}

	CreateTempleteResult struct {
		Templete *mysql_entity.TempleteModel
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
		Templete         *mysql_entity.TempleteModel
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
		Templete         *mysql_entity.TempleteModel
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
		Templetes []*mysql_entity.TempleteModel
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
		Templete         *mysql_entity.TempleteModel
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
