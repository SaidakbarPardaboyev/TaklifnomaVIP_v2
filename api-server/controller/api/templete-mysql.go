package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"saidakbar.origin/api-server/contracts"
	middlewares "saidakbar.origin/api-server/middleware"
	requestmodels "saidakbar.origin/api-server/request-models"
	templete_mysql "saidakbar.origin/services/templete-mysql"
)

type TemplateMysqlController interface {
	GetAllTempletes(ginContext *gin.Context)
	GetTempleteCount(ginContext *gin.Context)
	GetTempleteByID(ginContext *gin.Context)
	CreateTemplete(ginContext *gin.Context)
	UpdateTemplete(ginContext *gin.Context)
	DeleteTemplete(ginContext *gin.Context)
}

type templatMysqlController struct {
	templeteService templete_mysql.Service
}

func NewTemplateMysqlController(
	templeteService templete_mysql.Service,
) TemplateMysqlController {
	return &templatMysqlController{
		templeteService: templeteService,
	}
}

func (c *templatMysqlController) GetAllTempletes(ginContext *gin.Context) {
	var requestModel requestmodels.GetAllTempleteRequest

	if err := ginContext.ShouldBindQuery(&requestModel); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := middlewares.GetAccount(ginContext)

	switch result, failure := c.templeteService.GetAllTempletes(templete_mysql.GetAllTempletesModel{
		OrganizationID: account.ActiveOrganization.ID,
		Skip:           requestModel.Skip,
		Limit:          requestModel.Limit,
		Search:         requestModel.Search,
		StartDate:      requestModel.GetDateTimeStart(),
		EndDate:        requestModel.GetDateTimeEnd(),
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	default:
		viewModel := []contracts.TemplateMysqlContract{}
		for _, t := range result.Templetes {
			viewModel = append(viewModel, contracts.CreateTemplateMysqlContract(t))
		}
		ginContext.JSON(http.StatusOK, viewModel)
	}
}

func (c *templatMysqlController) GetTempleteCount(ginContext *gin.Context) {
	var requestModel requestmodels.GetTempleteCountRequest

	if err := ginContext.ShouldBindQuery(&requestModel); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := middlewares.GetAccount(ginContext)

	switch result, failure := c.templeteService.GetTempleteCount(templete_mysql.GetTempleteCountModel{
		OrganizationID: account.ActiveOrganization.ID,
		Search:         requestModel.Search,
		StartDate:      requestModel.GetDateTimeStart(),
		EndDate:        requestModel.GetDateTimeEnd(),
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	default:
		ginContext.JSON(http.StatusOK, result.Count)
	}
}

func (c *templatMysqlController) GetTempleteByID(ginContext *gin.Context) {
	var id string
	if id = ginContext.Param("id"); id == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	account := middlewares.GetAccount(ginContext)

	switch result, failure := c.templeteService.GetTempleteByID(templete_mysql.GetTempleteByIDModel{
		ID:             id,
		OrganizationID: account.ActiveOrganization.ID,
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	case result.TempleteNotFound:
		ginContext.JSON(http.StatusConflict, gin.H{"templete_not_found": result.TempleteNotFound})
	default:
		ginContext.JSON(http.StatusOK, contracts.CreateTemplateMysqlContract(result.Templete))
	}
}

func (c *templatMysqlController) CreateTemplete(ginContext *gin.Context) {
	var requestModel requestmodels.CreateTempleteRequest
	if err := ginContext.ShouldBindJSON(&requestModel); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := requestModel.Validate(); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := middlewares.GetAccount(ginContext)

	switch result, failure := c.templeteService.CreateTemplete(templete_mysql.CreateTempleteModel{
		OrganizationID: account.ActiveOrganization.ID,
		Name:           requestModel.Name,
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	default:
		ginContext.JSON(http.StatusOK, contracts.CreateTemplateMysqlContract(result.Templete))
	}
}

func (c *templatMysqlController) UpdateTemplete(ginContext *gin.Context) {
	var id string
	if id = ginContext.Param("id"); id == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var requestModel requestmodels.UpdateTempleteRequest
	if err := ginContext.ShouldBindJSON(&requestModel); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := requestModel.Validate(); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := middlewares.GetAccount(ginContext)

	switch result, failure := c.templeteService.UpdateTemplete(templete_mysql.UpdateTempleteModel{
		ID:             id,
		OrganizationID: account.ActiveOrganization.ID,
		Name:           requestModel.Name,
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	case result.TempleteNotFound:
		ginContext.JSON(http.StatusConflict, gin.H{"templete_not_found": result.TempleteNotFound})
	default:
		ginContext.JSON(http.StatusOK, contracts.CreateTemplateMysqlContract(result.Templete))
	}
}

func (c *templatMysqlController) DeleteTemplete(ginContext *gin.Context) {
	var id string
	if id = ginContext.Param("id"); id == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	account := middlewares.GetAccount(ginContext)

	switch result, failure := c.templeteService.DeleteTemplete(templete_mysql.DeleteTempleteModel{
		ID:             id,
		OrganizationID: account.ActiveOrganization.ID,
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	case result.TempleteNotFound:
		ginContext.JSON(http.StatusConflict, gin.H{"templete_not_found": result.TempleteNotFound})
	default:
		ginContext.JSON(http.StatusOK, contracts.CreateTemplateMysqlContract(result.Templete))
	}
}
