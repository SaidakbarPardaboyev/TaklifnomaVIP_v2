package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"saidakbar.origin/api-server/contracts"
	middlewares "saidakbar.origin/api-server/middleware"
	requestmodels "saidakbar.origin/api-server/request-models"
	"saidakbar.origin/core"
	payment_service "saidakbar.origin/services/payment"
)

type PaymentController interface {
	GeneratePaymeLink(ctx *gin.Context)
	PaymeWebhook(ctx *gin.Context)
}

type paymentController struct {
	paymentService payment_service.Service
}

func NewPaymentController(paymentService payment_service.Service) PaymentController {
	return &paymentController{paymentService: paymentService}
}

func (c *paymentController) GeneratePaymeLink(ctx *gin.Context) {
	var req requestmodels.GeneratePaymeLinkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc := middlewares.GetAccount(ctx)

	switch result, err := c.paymentService.GeneratePaymeLink(payment_service.GeneratePaymeLinkModel{
		AccountID: acc.ID,
		OrderID:   req.OrderID,
		Amount:    req.Amount,
	}); {
	case err != nil:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case result.OrderNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
	case result.Unauthorized:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
	case result.IncorrectAmount:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect amount"})
	default:
		ctx.JSON(http.StatusOK, contracts.CreatePaymentLinkContract(*result))
	}
}

func (c *paymentController) PaymeWebhook(ctx *gin.Context) {
	var req requestmodels.PaymeWebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}

	switch req.Method {
	case core.CheckPerformTransaction:
		c.checkPerformTransaction(ctx, req.Params)
	case core.CreateTransaction:
		c.createTransaction(ctx, req.Params)
	case core.PerformTransaction:
		c.performTransaction(ctx, req.Params)
	case core.CancelTransaction:
		c.cancelTransaction(ctx, req.Params)
	case core.CheckTransaction:
		c.checkTransaction(ctx, req.Params)
	case core.GetStatement:
		c.getStatement(ctx, req.Params)
	default:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, "unknown method", "unknown method", "unknown method"))
	}
}

func (c *paymentController) checkPerformTransaction(ctx *gin.Context, params json.RawMessage) {
	var req requestmodels.CheckPerformTransactionRequest
	if err := json.Unmarshal(params, &req); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}

	switch result, err := c.paymentService.CheckPerformTransaction(payment_service.CheckPerformTransactionModel{
		AccountID: req.Account.AccountID,
		OrderID:   req.Account.OrderID,
		Amount:    req.Amount,
	}); {
	case err != nil:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeInternalServerError, err.Error(), err.Error(), err.Error()))
	case !result.Succeed:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(result.ErrorCode, result.MessageUz, result.MessageRu, result.MessageEn))
	default:
		ctx.JSON(http.StatusOK, contracts.CreateCheckPerformTransactionContract(result.Order))
	}
}

func (c *paymentController) createTransaction(ctx *gin.Context, params json.RawMessage) {
	var req requestmodels.CreateTransactionRequest
	if err := json.Unmarshal(params, &req); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}

	switch result, err := c.paymentService.CreateTransaction(payment_service.CreateTransactionModel{
		AccountID: req.Account.AccountID,
		OrderID:   req.Account.OrderID,
		PaymentID: req.PaymentID,
		Amount:    req.Amount,
		TimePayme: req.TimePayme,
	}); {
	case err != nil:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeInternalServerError, err.Error(), err.Error(), err.Error()))
	case !result.Succeed:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(result.ErrorCode, result.MessageUz, result.MessageRu, result.MessageEn))
	default:
		ctx.JSON(http.StatusOK, contracts.BuildCreateTransactionContract(result.Transaction))
	}
}

func (c *paymentController) performTransaction(ctx *gin.Context, params json.RawMessage) {
	var req requestmodels.PerformTransactionRequest
	if err := json.Unmarshal(params, &req); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}

	switch result, err := c.paymentService.PerformTransaction(payment_service.PerformTransactionModel{
		PaymentID: req.ID,
	}); {
	case err != nil:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeInternalServerError, err.Error(), err.Error(), err.Error()))
	case !result.Succeed:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(result.ErrorCode, result.MessageUz, result.MessageRu, result.MessageEn))
	default:
		ctx.JSON(http.StatusOK, contracts.BuildPerformTransactionContract(result.Transaction))
	}
}

func (c *paymentController) cancelTransaction(ctx *gin.Context, params json.RawMessage) {
	var req requestmodels.CancelTransactionRequest
	if err := json.Unmarshal(params, &req); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}

	switch result, err := c.paymentService.CancelTransaction(payment_service.CancelTransactionModel{
		PaymentID: req.ID,
		Reason:    req.Reason,
	}); {
	case err != nil:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeInternalServerError, err.Error(), err.Error(), err.Error()))
	case !result.Succeed:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(result.ErrorCode, result.MessageUz, result.MessageRu, result.MessageEn))
	default:
		ctx.JSON(http.StatusOK, contracts.BuildCancelTransactionContract(result.Transaction))
	}
}

func (c *paymentController) checkTransaction(ctx *gin.Context, params json.RawMessage) {
	var req requestmodels.CheckTransactionRequest
	if err := json.Unmarshal(params, &req); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}

	switch result, err := c.paymentService.CheckTransaction(payment_service.CheckTransactionModel{
		PaymentID: req.ID,
	}); {
	case err != nil:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeInternalServerError, err.Error(), err.Error(), err.Error()))
	case !result.Succeed:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(result.ErrorCode, result.MessageUz, result.MessageRu, result.MessageEn))
	default:
		ctx.JSON(http.StatusOK, contracts.BuildCheckTransactionContract(result.Transaction))
	}
}

func (c *paymentController) getStatement(ctx *gin.Context, params json.RawMessage) {
	var req requestmodels.GetStatementRequest
	if err := json.Unmarshal(params, &req); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeBadRequest, err.Error(), err.Error(), err.Error()))
		return
	}

	switch result, err := c.paymentService.GetStatement(payment_service.GetStatementModel{
		From: req.From,
		To:   req.To,
	}); {
	case err != nil:
		ctx.JSON(http.StatusOK, contracts.CreatePaymeErrorContract(core.PaymeInternalServerError, err.Error(), err.Error(), err.Error()))
	default:
		ctx.JSON(http.StatusOK, contracts.BuildGetStatementContract(result.Transactions))
	}
}
