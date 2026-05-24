package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Cors(allowOrigins []string) func(ginContext *gin.Context) {
	corsConfig := cors.Config{}

	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = allowOrigins
	corsConfig.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodOptions}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	corsConfig.MaxAge = 24 * time.Hour

	return cors.New(corsConfig)
}
