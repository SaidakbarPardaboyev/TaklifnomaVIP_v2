package requestmodels

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type InvitationInfoRequest struct {
	Title         string    `json:"title"`
	GroomFullname string    `json:"groom_fullname"`
	BrideFullname string    `json:"bride_fullname"`
	PartnersNames string    `json:"partners_names"`
	Description   *string   `json:"description"`
	EventDate     time.Time `json:"event_date"`
	EventTime     string    `json:"event_time"`
	Address       string    `json:"address"`
	Location      string    `json:"location"`
	VenueName     *string   `json:"venue_name"`
	MainImage     *string   `json:"main_image"`
	Images        []string  `json:"images"`
	Story         *string   `json:"story"`
	Music         *string   `json:"music"`
}

type CreateOrderRequest struct {
	TemplateCode   string                `json:"template_code"`
	InvitationInfo InvitationInfoRequest `json:"invitation_info"`
}

func (r CreateOrderRequest) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.TemplateCode, validation.Required),
	); err != nil {
		return err
	}
	return validation.ValidateStruct(&r.InvitationInfo,
		validation.Field(&r.InvitationInfo.Title, validation.Required),
		validation.Field(&r.InvitationInfo.GroomFullname, validation.Required),
		validation.Field(&r.InvitationInfo.BrideFullname, validation.Required),
		validation.Field(&r.InvitationInfo.Address, validation.Required),
		validation.Field(&r.InvitationInfo.Location, validation.Required),
		validation.Field(&r.InvitationInfo.EventTime, validation.Required),
	)
}

type UpdateOrderRequest struct {
	InvitationInfo InvitationInfoRequest `json:"invitation_info"`
}

func (r UpdateOrderRequest) Validate() error {
	return validation.ValidateStruct(&r.InvitationInfo,
		validation.Field(&r.InvitationInfo.Title, validation.Required),
		validation.Field(&r.InvitationInfo.Address, validation.Required),
		validation.Field(&r.InvitationInfo.Location, validation.Required),
		validation.Field(&r.InvitationInfo.EventTime, validation.Required),
	)
}

type GetAllOrdersRequest struct {
	Status       *string `form:"status"`
	TemplateCode *string `form:"template_code"`
	FromDate     *string `form:"from_date"`
	ToDate       *string `form:"to_date"`
	Page         *int    `form:"page"`
	Limit        *int    `form:"limit"`
	SortBy       *string `form:"sort_by"`
	Order        *string `form:"order"`
}
