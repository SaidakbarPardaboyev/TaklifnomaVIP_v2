package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"saidakbar.origin/api-server/contracts"
	middlewares "saidakbar.origin/api-server/middleware"
	requestmodels "saidakbar.origin/api-server/request-models"
	order_service "saidakbar.origin/services/order"
)

type OrderController interface {
	Create(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type orderController struct {
	orderService order_service.Service
}

func NewOrderController(orderService order_service.Service) OrderController {
	return &orderController{orderService: orderService}
}

func (c *orderController) Create(ctx *gin.Context) {
	var req requestmodels.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc := middlewares.GetAccount(ctx)
	result, err := c.orderService.Create(order_service.CreateOrderModel{
		AccountID:    acc.ActiveOrganization.ID,
		TemplateCode: req.TemplateCode,
		InvitationInfo: order_service.InvitationInfoModel{
			Title:         req.InvitationInfo.Title,
			GroomFullname: req.InvitationInfo.GroomFullname,
			BrideFullname: req.InvitationInfo.BrideFullname,
			PartnersNames: req.InvitationInfo.PartnersNames,
			Description:   req.InvitationInfo.Description,
			EventDate:     req.InvitationInfo.EventDate,
			EventTime:     req.InvitationInfo.EventTime,
			Address:       req.InvitationInfo.Address,
			Location:      req.InvitationInfo.Location,
			VenueName:     req.InvitationInfo.VenueName,
			MainImage:     req.InvitationInfo.MainImage,
			Images:        req.InvitationInfo.Images,
			Story:         req.InvitationInfo.Story,
			Music:         req.InvitationInfo.Music,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, contracts.CreateOrderContract(result.Order))
}

func (c *orderController) GetAll(ctx *gin.Context) {
	var req requestmodels.GetAllOrdersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc := middlewares.GetAccount(ctx)

	model := order_service.GetAllOrdersModel{
		AccountID:    acc.ActiveOrganization.ID,
		Status:       req.Status,
		TemplateCode: req.TemplateCode,
		Page:         req.Page,
		Limit:        req.Limit,
		SortBy:       req.SortBy,
		Order:        req.Order,
	}

	if req.FromDate != nil {
		t, err := time.Parse("2006-01-02", *req.FromDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid from_date format, use YYYY-MM-DD"})
			return
		}
		model.FromDate = &t
	}
	if req.ToDate != nil {
		t, err := time.Parse("2006-01-02", *req.ToDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid to_date format, use YYYY-MM-DD"})
			return
		}
		model.ToDate = &t
	}

	result, err := c.orderService.GetAll(model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, contracts.CreateOrderListContract(result.Orders))
}

func (c *orderController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	acc := middlewares.GetAccount(ctx)

	switch result, err := c.orderService.GetByID(order_service.GetOrderByIDModel{
		ID:        id,
		AccountID: acc.ActiveOrganization.ID,
	}); {
	case err != nil:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case result.NotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
	case result.Unauthorized:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
	default:
		ctx.JSON(http.StatusOK, contracts.CreateOrderContract(result.Order))
	}
}

func (c *orderController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req requestmodels.UpdateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc := middlewares.GetAccount(ctx)

	switch result, err := c.orderService.Update(order_service.UpdateOrderModel{
		ID:        id,
		AccountID: acc.ActiveOrganization.ID,
		InvitationInfo: order_service.InvitationInfoModel{
			Title:         req.InvitationInfo.Title,
			GroomFullname: req.InvitationInfo.GroomFullname,
			BrideFullname: req.InvitationInfo.BrideFullname,
			PartnersNames: req.InvitationInfo.PartnersNames,
			Description:   req.InvitationInfo.Description,
			EventDate:     req.InvitationInfo.EventDate,
			EventTime:     req.InvitationInfo.EventTime,
			Address:       req.InvitationInfo.Address,
			Location:      req.InvitationInfo.Location,
			VenueName:     req.InvitationInfo.VenueName,
			MainImage:     req.InvitationInfo.MainImage,
			Images:        req.InvitationInfo.Images,
			Story:         req.InvitationInfo.Story,
			Music:         req.InvitationInfo.Music,
		},
	}); {
	case err != nil:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case result.NotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
	case result.Unauthorized:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
	case result.NotDraft:
		ctx.JSON(http.StatusConflict, gin.H{"error": "only draft orders can be updated"})
	default:
		ctx.JSON(http.StatusOK, contracts.CreateOrderContract(result.Order))
	}
}

func (c *orderController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	acc := middlewares.GetAccount(ctx)

	switch result, err := c.orderService.Delete(order_service.DeleteOrderModel{
		ID:        id,
		AccountID: acc.ActiveOrganization.ID,
	}); {
	case err != nil:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case result.NotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
	case result.Unauthorized:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
	case result.NotDraft:
		ctx.JSON(http.StatusConflict, gin.H{"error": "only draft orders can be cancelled"})
	default:
		ctx.JSON(http.StatusNoContent, nil)
	}
}
