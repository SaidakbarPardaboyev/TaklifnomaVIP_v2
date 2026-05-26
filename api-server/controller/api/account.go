package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"saidakbar.origin/api-server/contracts"
	middlewares "saidakbar.origin/api-server/middleware"
	requestmodels "saidakbar.origin/api-server/request-models"
	"saidakbar.origin/services/account"
)

type AccountController interface {
	VerifyCode(ginContext *gin.Context)
	GetMe(ginContext *gin.Context)
	UpdateMe(ginContext *gin.Context)
}

type accountController struct {
	accountService account.Service
}

func NewAccountController(accountService account.Service) AccountController {
	return &accountController{accountService: accountService}
}

func (c *accountController) VerifyCode(ginContext *gin.Context) {
	var req requestmodels.VerifyCodeRequest

	if err := ginContext.ShouldBindJSON(&req); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch result, failure := c.accountService.VerifyCode(account.VerifyCodeModel{
		Phone: req.Phone,
		Code:  req.Code,
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	case result.UserNotFound:
		ginContext.JSON(http.StatusNotFound, gin.H{"error": "phone not registered — open the bot first"})
	case result.InvalidCode:
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired code"})
	default:
		ginContext.JSON(http.StatusOK, gin.H{"token": result.Token})
	}
}

func (c *accountController) GetMe(ginContext *gin.Context) {
	acc := middlewares.GetAccount(ginContext)
	ginContext.JSON(http.StatusOK, contracts.CreateAccountContract(acc))
}

func (c *accountController) UpdateMe(ginContext *gin.Context) {
	var req requestmodels.UpdateAccountRequest
	if err := ginContext.ShouldBindJSON(&req); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc := middlewares.GetAccount(ginContext)

	switch result, err := c.accountService.UpdateAccount(account.UpdateAccountModel{
		ID:       acc.ID,
		FullName: req.FullName,
	}); {
	case err != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		ginContext.JSON(http.StatusOK, contracts.CreateAccountContractFromEntity(result.Account))
	}
}
