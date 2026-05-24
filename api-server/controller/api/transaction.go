package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"saidakbar.origin/api-server/contracts"
	middlewares "saidakbar.origin/api-server/middleware"
	transaction_service "saidakbar.origin/services/transaction"
)

type TransactionController interface {
	GetList(ctx *gin.Context)
	GetByOrderID(ctx *gin.Context)
}

type transactionController struct {
	transactionService transaction_service.Service
}

func NewTransactionController(transactionService transaction_service.Service) TransactionController {
	return &transactionController{transactionService: transactionService}
}

func (c *transactionController) GetList(ctx *gin.Context) {
	acc := middlewares.GetAccount(ctx)

	result, err := c.transactionService.GetList(transaction_service.GetListModel{
		AccountID: acc.ActiveOrganization.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, contracts.BuildTransactionListContract(result.Transactions))
}

func (c *transactionController) GetByOrderID(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	acc := middlewares.GetAccount(ctx)

	switch result, err := c.transactionService.GetByOrderID(transaction_service.GetByOrderIDModel{
		OrderID:   orderID,
		AccountID: acc.ActiveOrganization.ID,
	}); {
	case err != nil:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case result.NotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
	case result.Unauthorized:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
	default:
		ctx.JSON(http.StatusOK, contracts.BuildTransactionContract(result.Transaction))
	}
}
