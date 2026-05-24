package contracts

import (
	"time"

	"saidakbar.origin/db/mongo/entity"
)

type InvitationInfoContract struct {
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

type OrderContract struct {
	ID             string                  `json:"id"`
	AccountID      string                  `json:"account_id"`
	TemplateCode   string                  `json:"template_code"`
	InvitationInfo *InvitationInfoContract `json:"invitation_info"`
	Price          float64                 `json:"price"`
	Status         string                  `json:"status"`
	ViewCount      int64                   `json:"view_count"`
	ExpiresAt      *time.Time              `json:"expires_at"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

func CreateOrderContract(order *entity.OrderEntity) OrderContract {
	c := OrderContract{
		ID:           order.ID,
		AccountID:    order.AccountID,
		TemplateCode: order.TemplateCode,
		Price:        order.Price,
		Status:       order.Status,
		ViewCount:    order.ViewCount,
		ExpiresAt:    order.ExpiresAt,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}
	if order.InvitationInfo != nil {
		info := order.InvitationInfo
		c.InvitationInfo = &InvitationInfoContract{
			Title:         info.Title,
			GroomFullname: info.GroomFullname,
			BrideFullname: info.BrideFullname,
			PartnersNames: info.PartnersNames,
			Description:   info.Description,
			EventDate:     info.EventDate,
			EventTime:     info.EventTime,
			Address:       info.Address,
			Location:      info.Location,
			VenueName:     info.VenueName,
			MainImage:     info.MainImage,
			Images:        info.Images,
			Story:         info.Story,
			Music:         info.Music,
		}
	}
	return c
}

func CreateOrderListContract(orders []*entity.OrderEntity) []OrderContract {
	result := make([]OrderContract, 0, len(orders))
	for _, o := range orders {
		result = append(result, CreateOrderContract(o))
	}
	return result
}
