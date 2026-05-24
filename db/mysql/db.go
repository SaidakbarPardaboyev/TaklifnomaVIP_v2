package mysql

import (
	"fmt"

	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/plugins"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

func (d *database) GetDB() *gorm.DB {
	return d.db
}

func (d *database) Migrate() error {
	return d.db.AutoMigrate(
		mysql_entity.TempleteModel{},
		mysql_entity.AccountModel{},
		// TODO: add your GORM model structs here
	)
}

type Database interface {
	Migrate() error
	GetDB() *gorm.DB
}

func New(username, password, addr, databaseName string) (Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, addr, databaseName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NowFunc:                plugins.GetNow,
	})
	if err != nil {
		return nil, err
	}

	return &database{db: db}, nil
}
