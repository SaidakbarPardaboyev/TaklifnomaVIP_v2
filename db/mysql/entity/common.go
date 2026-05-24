package mysql_entity

import "time"

// BaseEntity is embedded in every MySQL GORM model.
type BaseEntity struct {
	ID        string     `gorm:"column:id;primaryKey;size:36"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	IsDeleted bool       `gorm:"column:is_deleted"`
}
