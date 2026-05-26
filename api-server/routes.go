package apiserver

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"saidakbar.origin/api-server/controller/api"
	middlewares "saidakbar.origin/api-server/middleware"
)

func (apiServer *apiServer) registerRoutes() {
	url := ginSwagger.URL("swagger/doc.json")
	apiServer.server.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// auth api
	var accountController = api.NewAccountController(apiServer.accountService)
	{
		apiServer.server.POST("/account/verify-code", accountController.VerifyCode)
	}

	// public api
	publicGroup := apiServer.server.Group("/public")
	{
		var invitationController = api.NewInvitationController(apiServer.orderService)
		publicGroup.GET("/i/:id", invitationController.GetPublic)
	}

	// payme api
	var paymentController = api.NewPaymentController(apiServer.paymentService)
	{
		apiServer.server.POST("/payme",
			middlewares.PaymeAuth(apiServer.paymeApiKeyProd, apiServer.paymeApiKeyStaging),
			paymentController.PaymeWebhook,
		)
	}

	// api
	apiGroup := apiServer.server.Group("/api", middlewares.AuthAccount(apiServer.tokenCache))
	{
		apiGroup.GET("/account/get-me", accountController.GetMe)
		apiGroup.PUT("/account/update-me", accountController.UpdateMe)

		var orderController = api.NewOrderController(apiServer.orderService)
		apiGroup.POST("/order/create", orderController.Create)
		apiGroup.GET("/order/list", orderController.GetAll)
		apiGroup.GET("/order/get/:id", orderController.GetByID)
		apiGroup.PUT("/order/update/:id", orderController.Update)
		apiGroup.DELETE("/order/delete/:id", orderController.Delete)

		apiGroup.POST("/payment/generate-link", paymentController.GeneratePaymeLink)

		var txController = api.NewTransactionController(apiServer.transactionService)
		apiGroup.GET("/transaction/list", txController.GetList)
		apiGroup.GET("/transaction/get/:order_id", txController.GetByOrderID)
	}
}
