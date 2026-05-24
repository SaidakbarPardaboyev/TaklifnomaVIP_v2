package middlewares

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"saidakbar.origin/api-server/contracts"
	"saidakbar.origin/core"
)

func PaymeAuth(apiKeyProd, apiKeyStaging string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			ctx.AbortWithStatusJSON(http.StatusOK,
				contracts.CreatePaymeErrorContract(core.PaymeAccessDenied, "Unauthorized", "Unauthorized", "Unauthorized"))
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusOK,
				contracts.CreatePaymeErrorContract(core.PaymeAccessDenied, "Unauthorized", "Unauthorized", "Unauthorized"))
			return
		}

		key := string(decoded)
		if key != core.PaymeApiKeyPrefix+apiKeyProd && key != core.PaymeApiKeyPrefix+apiKeyStaging {
			ctx.AbortWithStatusJSON(http.StatusOK,
				contracts.CreatePaymeErrorContract(core.PaymeAccessDenied, "Unauthorized", "Unauthorized", "Unauthorized"))
			return
		}

		ctx.Next()
	}
}
