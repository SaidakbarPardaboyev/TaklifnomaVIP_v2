package requestmodels

import validation "github.com/go-ozzo/ozzo-validation/v4"

type UpdateAccountRequest struct {
	FullName string `json:"full_name"`
}

func (r UpdateAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FullName, validation.Required, validation.Length(1, 255)),
	)
}
