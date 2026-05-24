package apiserver

import (
	"github.com/gin-gonic/gin"
	middlewares "saidakbar.origin/api-server/middleware"
	"saidakbar.origin/caching"
	_ "saidakbar.origin/docs"
	account_service "saidakbar.origin/services/account"
	templete_mongo "saidakbar.origin/services/templete-mongo"
	templete_mysql "saidakbar.origin/services/templete-mysql"
)

// @title           Templete API
// @version         1.0
// @description     Templete project API
// @host            localhost:8080
// @BasePath        /
// @schemes         http https
// @securityDefinitions.apiKey BearerAuth
// @in              header
// @name            Authorization

type Server interface {
	Start() error
}

type apiServer struct {
	bindAddr               string
	server                 *gin.Engine
	tokenCache             caching.TokenCache
	clientCredentialsCache caching.ClientCredentialsCache
	accountService         account_service.Service
	templeteMongoService   templete_mongo.Service
	templatMysqlService    templete_mysql.Service
	// TODO: add more service fields here
}

func (apiServer *apiServer) Start() error {
	return apiServer.server.Run(apiServer.bindAddr)
}

func NewApiServer(
	bindAddr string,
	allowedOrigins []string,
	tokenCache caching.TokenCache,
	clientCredentialsCache caching.ClientCredentialsCache,
	accountService account_service.Service,
	templeteMongoService templete_mongo.Service,
	templatMysqlService templete_mysql.Service,
	// TODO: add more services here
) Server {
	as := &apiServer{
		bindAddr:               bindAddr,
		server:                 gin.Default(),
		tokenCache:             tokenCache,
		clientCredentialsCache: clientCredentialsCache,
		accountService:         accountService,
		templeteMongoService:   templeteMongoService,
		templatMysqlService:    templatMysqlService,
	}

	as.server.Use(middlewares.Cors(allowedOrigins))
	as.registerRoutes()

	return as
}
