package requestmodels

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// region Get All
type GetAllTempleteRequest struct {
	Skip      *int    `form:"skip"`
	Limit     *int    `form:"limit"`
	Search    *string `form:"search"`
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
}

func (rm *GetAllTempleteRequest) GetDateTimeStart() *time.Time {
	return parseDateTime(rm.StartDate)
}

func (rm *GetAllTempleteRequest) GetDateTimeEnd() *time.Time {
	return parseDateTime(rm.EndDate)
}

// region Get Count
type GetTempleteCountRequest struct {
	Search    *string `form:"search"`
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
}

func (rm *GetTempleteCountRequest) GetDateTimeStart() *time.Time {
	return parseDateTime(rm.StartDate)
}

func (rm *GetTempleteCountRequest) GetDateTimeEnd() *time.Time {
	return parseDateTime(rm.EndDate)
}

// region Create
type CreateTempleteRequest struct {
	Name string `json:"name"`
	// TODO: add domain-specific create fields
}

func (r CreateTempleteRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
	)
}

// region Update
type UpdateTempleteRequest struct {
	Name string `json:"name"`
	// TODO: add domain-specific update fields
}

func (r UpdateTempleteRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
	)
}
