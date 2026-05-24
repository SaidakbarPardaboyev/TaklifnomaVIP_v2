package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	mysql_db "saidakbar.origin/db/mysql"
	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/repository/models"
)

type TransactionRepository interface {
	Create(ctx context.Context, model *mysql_entity.TransactionModel, tx *gorm.DB) error
	GetByPaymentID(ctx context.Context, req *models.GetTransactionByPaymentIDRequest, tx *gorm.DB) (*mysql_entity.TransactionModel, error)
	GetByOrderID(ctx context.Context, req *models.GetTransactionByOrderIDRequest, tx *gorm.DB) (*mysql_entity.TransactionModel, error)
	GetList(ctx context.Context, req *models.GetListTransactionRequest, tx *gorm.DB) ([]*mysql_entity.TransactionModel, error)
	Update(ctx context.Context, model *mysql_entity.TransactionModel, tx *gorm.DB) error
	GetDB() *gorm.DB
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(mysqlDB mysql_db.Database) TransactionRepository {
	return &transactionRepository{db: mysqlDB.GetDB()}
}

func (r *transactionRepository) GetDB() *gorm.DB { return r.db }

func (r *transactionRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *transactionRepository) Create(ctx context.Context, model *mysql_entity.TransactionModel, tx *gorm.DB) error {
	return r.getDB(tx).WithContext(ctx).Create(model).Error
}

func (r *transactionRepository) GetByPaymentID(ctx context.Context, req *models.GetTransactionByPaymentIDRequest, tx *gorm.DB) (*mysql_entity.TransactionModel, error) {
	var result mysql_entity.TransactionModel
	err := r.getDB(tx).WithContext(ctx).
		Where(mysql_entity.FieldTransactionPaymentID+" = ?", req.PaymentID).
		Order(mysql_entity.FieldTransactionCreatedTime + " DESC").
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *transactionRepository) GetByOrderID(ctx context.Context, req *models.GetTransactionByOrderIDRequest, tx *gorm.DB) (*mysql_entity.TransactionModel, error) {
	var result mysql_entity.TransactionModel
	err := r.getDB(tx).WithContext(ctx).
		Where(mysql_entity.FieldTransactionOrderID+" = ?", req.OrderID).
		Order(mysql_entity.FieldTransactionCreatedTime + " DESC").
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *transactionRepository) GetList(ctx context.Context, req *models.GetListTransactionRequest, tx *gorm.DB) ([]*mysql_entity.TransactionModel, error) {
	db := r.getDB(tx).WithContext(ctx).Model(&mysql_entity.TransactionModel{})

	if req.AccountID != nil {
		db = db.Where(mysql_entity.FieldTransactionAccountID+" = ?", *req.AccountID)
	}
	if req.OrderID != nil {
		db = db.Where(mysql_entity.FieldTransactionOrderID+" = ?", *req.OrderID)
	}
	if req.PaymentIDs != nil && len(*req.PaymentIDs) > 0 {
		db = db.Where(mysql_entity.FieldTransactionPaymentID+" IN ?", *req.PaymentIDs)
	}
	if req.States != nil && len(*req.States) > 0 {
		db = db.Where(mysql_entity.FieldTransactionState+" IN ?", *req.States)
	}
	if req.FromDate != nil {
		db = db.Where(mysql_entity.FieldTransactionCreatedTime+" >= ?", *req.FromDate)
	}
	if req.ToDate != nil {
		db = db.Where(mysql_entity.FieldTransactionCreatedTime+" <= ?", *req.ToDate)
	}
	if req.Limit != nil {
		db = db.Limit(*req.Limit)
		if req.Page != nil && *req.Page > 1 {
			db = db.Offset(*req.Limit * (*req.Page - 1))
		}
	}

	var results []*mysql_entity.TransactionModel
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *transactionRepository) Update(ctx context.Context, model *mysql_entity.TransactionModel, tx *gorm.DB) error {
	return r.getDB(tx).WithContext(ctx).Save(model).Error
}

