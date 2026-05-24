package payment_service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"saidakbar.origin/config"
	"saidakbar.origin/core"
	"saidakbar.origin/db/mongo/entity"
	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/plugins"
	"saidakbar.origin/repository"
	"saidakbar.origin/repository/models"
)

type paymentService struct {
	cfg         *config.Config
	accountRepo repository.AccountRepository
	orderRepo   repository.OrderRepository
	txRepo      repository.TransactionRepository
}

func NewPaymentService(
	cfg *config.Config,
	accountRepo repository.AccountRepository,
	orderRepo repository.OrderRepository,
	txRepo repository.TransactionRepository,
) Service {
	return &paymentService{
		cfg:         cfg,
		accountRepo: accountRepo,
		orderRepo:   orderRepo,
		txRepo:      txRepo,
	}
}

func (s *paymentService) GeneratePaymeLink(model GeneratePaymeLinkModel) (*GeneratePaymeLinkResult, error) {
	result := new(GeneratePaymeLinkResult)

	order, err := s.orderRepo.GetByID(model.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		result.OrderNotFound = true
		return result, nil
	}
	if order.AccountID != model.AccountID {
		result.Unauthorized = true
		return result, nil
	}
	if plugins.ConvertToPaymeAmount(order.Price) != plugins.ConvertToPaymeAmount(model.Amount) {
		result.IncorrectAmount = true
		return result, nil
	}

	params := fmt.Sprintf("m=%s;ac.order_id=%s;ac.account_id=%s;a=%d;c=%s;ct=5000;",
		s.cfg.PaymeMerchantID,
		order.ID,
		model.AccountID,
		int(plugins.ConvertToPaymeAmount(model.Amount)),
		s.cfg.PaymeRedirectionLink,
	)
	encoded := base64.StdEncoding.EncodeToString([]byte(params))
	link := "https://checkout.paycom.uz/" + encoded

	result.Link = &link
	return result, nil
}

func (s *paymentService) CheckPerformTransaction(model CheckPerformTransactionModel) (*CheckPerformTransactionResult, error) {
	result := new(CheckPerformTransactionResult)

	account, err := s.accountRepo.GetByID(model.AccountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		result.AccountNotFound = true
		result.ErrorCode = core.PaymeErrorNotFound
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageUserNotFound, model.AccountID)
		return result, nil
	}

	order, err := s.orderRepo.GetByID(model.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		result.OrderNotFound = true
		result.ErrorCode = core.PaymeErrorNotFound
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageOrderNotFound, model.OrderID)
		return result, nil
	}
	if order.AccountID != account.ID {
		result.Unauthorized = true
		result.ErrorCode = core.PaymeErrorNotFound
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageOrderNotBelongToUser, model.AccountID, model.OrderID)
		return result, nil
	}
	if plugins.ConvertToPaymeAmount(order.Price) != model.Amount {
		result.IncorrectAmount = true
		result.ErrorCode = core.PaymeIncorrectAmount
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageIncorrectAmount, model.Amount, plugins.ConvertToPaymeAmount(order.Price))
		return result, nil
	}
	if order.Status != entity.OrderStatusDraft && order.Status != entity.OrderStatusCancelled {
		result.AlreadyPaid = true
		result.ErrorCode = core.PaymeBadRequest
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageOrderAlreadyPaid)
		return result, nil
	}

	result.Succeed = true
	result.Order = order
	return result, nil
}

func (s *paymentService) CreateTransaction(model CreateTransactionModel) (*CreateTransactionResult, error) {
	result := new(CreateTransactionResult)
	ctx := context.Background()
	db := s.txRepo.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	account, err := s.accountRepo.GetByID(model.AccountID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if account == nil {
		_ = tx.Rollback()
		result.AccountNotFound = true
		result.ErrorCode = core.PaymeErrorNotFound
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageUserNotFound, model.AccountID)
		return result, nil
	}

	order, err := s.orderRepo.GetByID(model.OrderID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if order == nil {
		_ = tx.Rollback()
		result.OrderNotFound = true
		result.ErrorCode = core.PaymeErrorNotFound
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageOrderNotFound, model.OrderID)
		return result, nil
	}
	if order.AccountID != account.ID {
		_ = tx.Rollback()
		result.Unauthorized = true
		result.ErrorCode = core.PaymeErrorNotFound
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageOrderNotBelongToUser, model.AccountID, model.OrderID)
		return result, nil
	}
	if plugins.ConvertToPaymeAmount(order.Price) != float64(model.Amount) {
		_ = tx.Rollback()
		result.IncorrectAmount = true
		result.ErrorCode = core.PaymeIncorrectAmount
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageIncorrectAmount, float64(model.Amount), plugins.ConvertToPaymeAmount(order.Price))
		return result, nil
	}
	if order.Status != entity.OrderStatusDraft && order.Status != entity.OrderStatusCancelled {
		_ = tx.Rollback()
		result.AlreadyPaid = true
		result.ErrorCode = core.PaymeBadRequest
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageOrderAlreadyPaid)
		return result, nil
	}

	// return existing transaction if payment_id already seen
	existing, err := s.txRepo.GetByPaymentID(ctx, &models.GetTransactionByPaymentIDRequest{PaymentID: model.PaymentID}, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if existing != nil {
		_ = tx.Commit()
		result.Transaction = existing
		result.Succeed = true
		return result, nil
	}

	// reject if active/completed transaction already exists for this order
	list, err := s.txRepo.GetList(ctx, &models.GetListTransactionRequest{
		OrderID: &model.OrderID,
		States:  &[]int{core.TransactionStateCreated, core.TransactionStateCompleted},
	}, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if len(list) > 0 {
		_ = tx.Rollback()
		result.TransactionExists = true
		result.ErrorCode = core.PaymeTransactionExists
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageTransactionExists)
		return result, nil
	}

	createdTime := plugins.ConverterFromPaymeTimeFormat(model.TimePayme)
	transaction := &mysql_entity.TransactionModel{
		ID:          uuid.NewString(),
		PaymentID:   model.PaymentID,
		Amount:      plugins.ConvertToPaymeAmount(order.Price),
		CreatedTime: &createdTime,
		State:       core.TransactionStateCreated,
		TimePayme:   model.TimePayme,
		AccountID:   account.ID,
		OrderID:     order.ID,
	}

	if err = s.txRepo.Create(ctx, transaction, tx); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	result.Transaction = transaction
	result.Succeed = true
	return result, nil
}

func (s *paymentService) PerformTransaction(model PerformTransactionModel) (*PerformTransactionResult, error) {
	result := new(PerformTransactionResult)
	ctx := context.Background()
	db := s.txRepo.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	transaction, err := s.txRepo.GetByPaymentID(ctx, &models.GetTransactionByPaymentIDRequest{PaymentID: model.PaymentID}, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if transaction == nil {
		_ = tx.Rollback()
		result.TransactionNotFound = true
		result.ErrorCode = core.PaymeErrorNotFound
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageTransactionNotFound, model.PaymentID)
		return result, nil
	}
	if transaction.State == core.TransactionStateCompleted {
		_ = tx.Commit()
		result.Transaction = transaction
		result.Succeed = true
		return result, nil
	}

	performTime := time.Now()
	transaction.PerformTime = &performTime
	transaction.State = core.TransactionStateCompleted

	if err = s.txRepo.Update(ctx, transaction, tx); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err = s.orderRepo.UpdateStatus(transaction.OrderID, entity.OrderStatusActive); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	result.Transaction = transaction
	result.Succeed = true
	return result, nil
}

func (s *paymentService) CancelTransaction(model CancelTransactionModel) (*CancelTransactionResult, error) {
	result := new(CancelTransactionResult)
	ctx := context.Background()
	db := s.txRepo.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	transaction, err := s.txRepo.GetByPaymentID(ctx, &models.GetTransactionByPaymentIDRequest{PaymentID: model.PaymentID}, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if transaction == nil {
		_ = tx.Rollback()
		result.TransactionNotFound = true
		result.ErrorCode = core.PaymeErrorNotFound
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageTransactionNotFound, model.PaymentID)
		return result, nil
	}
	if transaction.State == core.TransactionStateCanceledWhenCreated || transaction.State == core.TransactionStateCanceledWhenPerformed {
		_ = tx.Commit()
		result.Transaction = transaction
		result.Succeed = true
		return result, nil
	}

	cancelTime := time.Now()
	state := core.TransactionStateCanceledWhenPerformed
	reason := core.TransactionCancelReasonWhenPerformed
	if transaction.State == core.TransactionStateCreated {
		state = core.TransactionStateCanceledWhenCreated
		reason = core.TransactionCancelReasonWhenCreated
	}

	transaction.CancelTime = &cancelTime
	transaction.State = state
	transaction.Reason = &reason

	if err = s.txRepo.Update(ctx, transaction, tx); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err = s.orderRepo.UpdateStatus(transaction.OrderID, entity.OrderStatusCancelled); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	result.Transaction = transaction
	result.Succeed = true
	return result, nil
}

func (s *paymentService) CheckTransaction(model CheckTransactionModel) (*CheckTransactionResult, error) {
	result := new(CheckTransactionResult)
	ctx := context.Background()
	db := s.txRepo.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	transaction, err := s.txRepo.GetByPaymentID(ctx, &models.GetTransactionByPaymentIDRequest{PaymentID: model.PaymentID}, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if transaction == nil {
		_ = tx.Rollback()
		result.TransactionNotFound = true
		result.ErrorCode = core.PaymeErrorNotFound
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageTransactionNotFound, model.PaymentID)
		return result, nil
	}

	order, err := s.orderRepo.GetByID(transaction.OrderID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if order == nil {
		_ = tx.Rollback()
		result.OrderNotFound = true
		result.ErrorCode = core.PaymeErrorNotFound
		result.MessageUz, result.MessageRu, result.MessageEn = core.CreatePaymeMessage(core.PaymeErrorMessageOrderNotFound, transaction.OrderID)
		return result, nil
	}

	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	result.Transaction = transaction
	result.Succeed = true
	return result, nil
}

func (s *paymentService) GetStatement(model GetStatementModel) (*GetStatementResult, error) {
	from := plugins.ConverterFromPaymeTimeFormat(model.From)
	to := plugins.ConverterFromPaymeTimeFormat(model.To)

	txs, err := s.txRepo.GetList(context.Background(), &models.GetListTransactionRequest{
		FromDate: &from,
		ToDate:   &to,
	}, nil)
	if err != nil {
		return nil, err
	}
	return &GetStatementResult{Transactions: txs}, nil
}
