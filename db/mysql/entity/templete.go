package mysql_entity

import "gorm.io/gorm"

// TempleteModel is the example MySQL table.
// Rename this file and struct to match your domain (e.g. UserModel, OrderModel).
type TempleteModel struct {
	BaseEntity
	OrganizationID string `gorm:"column:organization_id;index"`
	Name           string `gorm:"column:name"`
	// TODO: add domain-specific columns here
}

func (TempleteModel) TableName() string {
	return "templete"
}

func (TempleteModel) FilterID(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldTempleteID+" = ?", id)
	}
}

func (TempleteModel) FilterOrganizationID(organizationID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldTempleteOrganizationID+" = ?", organizationID)
	}
}

func (TempleteModel) FilterIsDeleted(isDeleted bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldTempleteIsDeleted+" = ?", isDeleted)
	}
}

func (TempleteModel) FilterName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldTempleteName+" LIKE ?", "%"+name+"%")
	}
}

func (TempleteModel) FilterCreatedAtGte(t interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldTempleteCreatedAt+" >= ?", t)
	}
}

func (TempleteModel) FilterCreatedAtLte(t interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldTempleteCreatedAt+" <= ?", t)
	}
}

func (TempleteModel) OrderByCreatedAtDesc() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(FieldTempleteCreatedAt + " DESC")
	}
}

func (TempleteModel) OrderByCreatedAtAsc() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(FieldTempleteCreatedAt + " ASC")
	}
}
