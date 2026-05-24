package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"saidakbar.origin/caching"
	"saidakbar.origin/caching/models"
)

const (
	ContextUserKey   = "user"
	TokenTypeAccount = 0
	TokenTypeClient  = 1
)

func AuthAccount(client caching.TokenCache) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		isOk := false

		for key, value := range ginContext.Request.Header {
			if strings.ToLower(key) == "authorization" {
				user, err := client.GetAccount(value[0])
				if err == nil {
					if user.TokenType == TokenTypeAccount {
						ginContext.Set(ContextUserKey, user)
						isOk = true
						break
					}
				}
			}
		}

		if !isOk {
			ginContext.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func GetAccount(ginContext *gin.Context) *models.Account {
	var userObject, _ = ginContext.Get(ContextUserKey)
	var account, _ = userObject.(*models.Account)

	return account
}
