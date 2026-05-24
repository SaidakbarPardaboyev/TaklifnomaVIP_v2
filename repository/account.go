package repository

import (
	"errors"

	"gorm.io/gorm"
	mysql_db "saidakbar.origin/db/mysql"
	mysql_entity "saidakbar.origin/db/mysql/entity"
)

type AccountRepository interface {
	Create(account *mysql_entity.AccountModel) (err error)
	Update(account *mysql_entity.AccountModel) (err error)
	Activate(id string) (err error)
	GetByID(id string) (account *mysql_entity.AccountModel, err error)
	GetByPhone(phone string) (account *mysql_entity.AccountModel, err error)
	GetByChatID(chatID int64) (account *mysql_entity.AccountModel, err error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db mysql_db.Database) AccountRepository {
	return &accountRepository{db: db.GetDB()}
}

func (r *accountRepository) Create(account *mysql_entity.AccountModel) (err error) {
	return r.db.Create(account).Error
}

func (r *accountRepository) Update(account *mysql_entity.AccountModel) (err error) {
	return r.db.Model(account).Updates(map[string]interface{}{
		mysql_entity.FieldAccountFullName:  account.FullName,
		mysql_entity.FieldAccountChatID:    account.ChatID,
		mysql_entity.FieldAccountUpdatedAt: account.UpdatedAt,
	}).Error
}

func (r *accountRepository) Activate(id string) (err error) {
	return r.db.Model(&mysql_entity.AccountModel{}).
		Where(mysql_entity.FieldAccountID+" = ?", id).
		Update(mysql_entity.FieldAccountIsActive, true).Error
}

func (r *accountRepository) GetByID(id string) (account *mysql_entity.AccountModel, err error) {
	var entity mysql_entity.AccountModel
	err = r.db.Scopes(
		entity.FilterID(id),
		entity.FilterIsDeleted(false),
	).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}
	return
}

func (r *accountRepository) GetByPhone(phone string) (account *mysql_entity.AccountModel, err error) {
	var entity mysql_entity.AccountModel
	err = r.db.Scopes(
		entity.FilterPhone(phone),
		entity.FilterIsDeleted(false),
	).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}
	return
}

func (r *accountRepository) GetByChatID(chatID int64) (account *mysql_entity.AccountModel, err error) {
	var entity mysql_entity.AccountModel
	err = r.db.Scopes(
		entity.FilterChatID(chatID),
		entity.FilterIsDeleted(false),
	).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}
	return
}
