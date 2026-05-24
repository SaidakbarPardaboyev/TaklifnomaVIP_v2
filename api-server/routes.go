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

	var accountController = api.NewAccountController(apiServer.accountService)

	// account (no auth middleware)
	apiServer.server.POST("/account/verify-code", accountController.VerifyCode)

	// api
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
}
