package transaction_service

import (
	"context"

	"saidakbar.origin/repository"
	"saidakbar.origin/repository/models"
)

type transactionService struct {
	txRepo    repository.TransactionRepository
	orderRepo repository.OrderRepository
}

func NewTransactionService(txRepo repository.TransactionRepository, orderRepo repository.OrderRepository) Service {
	return &transactionService{txRepo: txRepo, orderRepo: orderRepo}
}

func (s *transactionService) GetByOrderID(model GetByOrderIDModel) (*GetByOrderIDResult, error) {
	result := new(GetByOrderIDResult)

	order, err := s.orderRepo.GetByID(model.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		result.NotFound = true
		return result, nil
	}
	if order.AccountID != model.AccountID {
		result.Unauthorized = true
		return result, nil
	}

	tx, err := s.txRepo.GetByOrderID(context.Background(), &models.GetTransactionByOrderIDRequest{OrderID: model.OrderID}, nil)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		result.NotFound = true
		return result, nil
	}

	result.Transaction = tx
	return result, nil
}

func (s *transactionService) GetList(model GetListModel) (*GetListResult, error) {
	result := new(GetListResult)

	limit := 10
	if model.Limit != nil && *model.Limit > 0 && *model.Limit <= 100 {
		limit = *model.Limit
	}

	txs, err := s.txRepo.GetList(context.Background(), &models.GetListTransactionRequest{
		AccountID: &model.AccountID,
		Limit:     &limit,
		Page:      model.Page,
	}, nil)
	if err != nil {
		return nil, err
	}

	result.Transactions = txs
	return result, nil
}
