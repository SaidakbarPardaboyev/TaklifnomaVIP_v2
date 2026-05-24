package entity

import "time"

const (
	OrderStatusDraft     = "draft"
	OrderStatusActive    = "active"
	OrderStatusExpired   = "expired"
	OrderStatusCancelled = "cancelled"
)

type InvitationInfo struct {
	Title         string    `bson:"title"          json:"title"`
	GroomFullname string    `bson:"groom_fullname" json:"groom_fullname"`
	BrideFullname string    `bson:"bride_fullname" json:"bride_fullname"`
	PartnersNames string    `bson:"partners_names" json:"partners_names"`
	Description   *string   `bson:"description"    json:"description"`
	EventDate     time.Time `bson:"event_date"     json:"event_date"`
	EventTime     string    `bson:"event_time"     json:"event_time"`
	Address       string    `bson:"address"        json:"address"`
	Location      string    `bson:"location"       json:"location"`
	VenueName     *string   `bson:"venue_name"     json:"venue_name"`
	MainImage     *string   `bson:"main_image"     json:"main_image"`
	Images        []string  `bson:"images"         json:"images"`
	Story         *string   `bson:"story"          json:"story"`
	Music         *string   `bson:"music"          json:"music"`
}

type OrderEntity struct {
	ID             string          `bson:"_id"             json:"id"`
	AccountID      string          `bson:"account_id"      json:"account_id"`
	TemplateCode   string          `bson:"template_code"   json:"template_code"`
	InvitationInfo *InvitationInfo `bson:"invitation_info" json:"invitation_info"`
	Price          float64         `bson:"price"           json:"price"`
	Status         string          `bson:"status"          json:"status"`
	ViewCount      int64           `bson:"view_count"      json:"view_count"`
	ExpiresAt      *time.Time      `bson:"expires_at"      json:"expires_at"`
	CreatedAt      time.Time       `bson:"created_at"      json:"created_at"`
	UpdatedAt      time.Time       `bson:"updated_at"      json:"updated_at"`
}
