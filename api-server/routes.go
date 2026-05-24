package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"saidakbar.origin/api-server/controller/api"
	middlewares "saidakbar.origin/api-server/middleware"
)

func (apiServer *apiServer) registerRoutes() {
	url := ginSwagger.URL("swagger/doc.json")
	apiServer.server.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	var accountController = api.NewAccountController(apiServer.accountService)

	// account (no auth middleware)
	apiServer.server.POST("/account/verify-code", accountController.VerifyCode)

	// public invitation view
	apiServer.server.GET("/i/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		_ = apiServer.orderService.IncrementViewCount(id)

		result, err := apiServer.orderService.GetPublic(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if result.NotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invitation not found"})
			return
		}
		ctx.JSON(http.StatusOK, result.Order)
	})

	// payme webhook (Basic Auth, no user token)
	var paymentController = api.NewPaymentController(apiServer.paymentService)
	apiServer.server.POST("/payme",
		middlewares.PaymeAuth(apiServer.paymeApiKeyProd, apiServer.paymeApiKeyStaging),
		paymentController.PaymeWebhook,
	)

	// api (user token required)
	apiGroup := apiServer.server.Group("/api", middlewares.AuthAccount(apiServer.tokenCache))
	{
		apiGroup.GET("/account/get-me", accountController.GetMe)

		var templeteMongoController = api.NewTempleteMongoController(apiServer.templeteMongoService)
		apiGroup.GET("/templete-mongo/get", templeteMongoController.GetAllTempletes)
		apiGroup.GET("/templete-mongo/count", templeteMongoController.GetTempleteCount)
		apiGroup.GET("/templete-mongo/get/:id", templeteMongoController.GetTempleteByID)
		apiGroup.POST("/templete-mongo/create", templeteMongoController.CreateTemplete)
		apiGroup.POST("/templete-mongo/update/:id", templeteMongoController.UpdateTemplete)
		apiGroup.POST("/templete-mongo/delete/:id", templeteMongoController.DeleteTemplete)

		var templatMysqlController = api.NewTemplateMysqlController(apiServer.templatMysqlService)
		apiGroup.GET("/templete-mysql/get", templatMysqlController.GetAllTempletes)
		apiGroup.GET("/templete-mysql/count", templatMysqlController.GetTempleteCount)
		apiGroup.GET("/templete-mysql/get/:id", templatMysqlController.GetTempleteByID)
		apiGroup.POST("/templete-mysql/create", templatMysqlController.CreateTemplete)
		apiGroup.POST("/templete-mysql/update/:id", templatMysqlController.UpdateTemplete)
		apiGroup.POST("/templete-mysql/delete/:id", templatMysqlController.DeleteTemplete)
	}

	// v1 (user token required)
	v1 := apiServer.server.Group("/v1", middlewares.AuthAccount(apiServer.tokenCache))
	{
		var orderController = api.NewOrderController(apiServer.orderService)
		v1.POST("/orders", orderController.Create)
		v1.GET("/orders", orderController.GetAll)
		v1.GET("/orders/:id", orderController.GetByID)
		v1.PUT("/orders/:id", orderController.Update)
		v1.DELETE("/orders/:id", orderController.Delete)

		v1.POST("/payment/generate-link", paymentController.GeneratePaymeLink)

		var txController = api.NewTransactionController(apiServer.transactionService)
		v1.GET("/transactions", txController.GetList)
		v1.GET("/transactions/:order_id", txController.GetByOrderID)
	}
}
