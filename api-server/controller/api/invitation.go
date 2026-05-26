package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"saidakbar.origin/api-server/contracts"
	order_service "saidakbar.origin/services/order"
)

type InvitationController interface {
	GetPublic(ctx *gin.Context)
}

type invitationController struct {
	orderService order_service.Service
}

func NewInvitationController(orderService order_service.Service) InvitationController {
	return &invitationController{orderService: orderService}
}

func (c *invitationController) GetPublic(ctx *gin.Context) {
	id := ctx.Param("id")
	_ = c.orderService.IncrementViewCount(id)

	switch result, err := c.orderService.GetPublic(id); {
	case err != nil:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case result.NotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"error": "invitation not found"})
	default:
		ctx.JSON(http.StatusOK, contracts.CreateOrderContract(result.Order))
	}
}
