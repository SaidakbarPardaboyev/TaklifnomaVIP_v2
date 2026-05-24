package requestmodels

import validation "github.com/go-ozzo/ozzo-validation/v4"

type VerifyCodeRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (r VerifyCodeRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Phone, validation.Required),
		validation.Field(&r.Code, validation.Required, validation.Length(6, 6)),
	)
}
