package models

import "time"

type Account struct {
	ID        string     `json:"id"`
	FullName  string     `json:"full_name"`
	Phone     string     `json:"phone"`
	ChatID    int64      `json:"chat_id"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	IsDeleted bool       `json:"is_deleted"`
	TokenType byte       `json:"token_type"`
}

type Organization struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Inn   *string `json:"inn"`
	Pinfl *string `json:"pinfl"`
}
