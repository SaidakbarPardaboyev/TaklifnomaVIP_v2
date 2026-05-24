package models

import "time"

type GetAllTempleteFilter struct {
	OrganizationID string
	Skip           *int
	Limit          *int
	Search         *string
	StartDate      *time.Time
	EndDate        *time.Time
}

type GetTempleteCountFilter struct {
	OrganizationID string
	Search         *string
	StartDate      *time.Time
	EndDate        *time.Time
}
