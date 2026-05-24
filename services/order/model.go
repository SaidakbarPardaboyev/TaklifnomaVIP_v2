package order_service

import (
	"time"

	"saidakbar.origin/db/mongo/entity"
)

type InvitationInfoModel struct {
	Title         string
	GroomFullname string
	BrideFullname string
	PartnersNames string
	Description   *string
	EventDate     time.Time
	EventTime     string
	Address       string
	Location      string
	VenueName     *string
	MainImage     *string
	Images        []string
	Story         *string
	Music         *string
}

type CreateOrderModel struct {
	AccountID      string
	TemplateCode   string
	InvitationInfo InvitationInfoModel
	Price          float64
}

type CreateOrderResult struct {
	Order *entity.OrderEntity
}

type GetOrderByIDModel struct {
	ID        string
	AccountID string
}

type GetOrderByIDResult struct {
	Order        *entity.OrderEntity
	NotFound     bool
	Unauthorized bool
}

type GetAllOrdersModel struct {
	AccountID    string
	Status       *string
	TemplateCode *string
	FromDate     *time.Time
	ToDate       *time.Time
	Page         *int
	Limit        *int
	SortBy       *string
	Order        *string
}

type GetAllOrdersResult struct {
	Orders []*entity.OrderEntity
}

type UpdateOrderModel struct {
	ID             string
	AccountID      string
	InvitationInfo InvitationInfoModel
}

type UpdateOrderResult struct {
	Order        *entity.OrderEntity
	NotFound     bool
	Unauthorized bool
	NotDraft     bool
}

type DeleteOrderModel struct {
	ID        string
	AccountID string
}

type DeleteOrderResult struct {
	NotFound     bool
	Unauthorized bool
	NotDraft     bool
}
