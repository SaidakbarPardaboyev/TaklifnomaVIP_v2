package mysql_entity

import "gorm.io/gorm"

type AccountModel struct {
	BaseEntity
	FullName string `gorm:"column:full_name;size:255"`
	Phone    string `gorm:"column:phone;uniqueIndex;size:20"`
	ChatID   int64  `gorm:"column:chat_id"`
	IsActive bool   `gorm:"column:is_active"`
}

func (AccountModel) TableName() string {
	return "accounts"
}

func (AccountModel) FilterID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldAccountID+" = ?", id)
	}
}

func (AccountModel) FilterPhone(phone string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldAccountPhone+" = ?", phone)
	}
}

func (AccountModel) FilterChatID(chatID int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldAccountChatID+" = ?", chatID)
	}
}

func (AccountModel) FilterIsDeleted(isDeleted bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldAccountIsDeleted+" = ?", isDeleted)
	}
}
