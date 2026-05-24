package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"saidakbar.origin/api-server/contracts"
	middlewares "saidakbar.origin/api-server/middleware"
	requestmodels "saidakbar.origin/api-server/request-models"
	templete_mongo "saidakbar.origin/services/templete-mongo"
)

type TempleteMongoController interface {
	GetAllTempletes(ginContext *gin.Context)
	GetTempleteCount(ginContext *gin.Context)
	GetTempleteByID(ginContext *gin.Context)
	CreateTemplete(ginContext *gin.Context)
	UpdateTemplete(ginContext *gin.Context)
	DeleteTemplete(ginContext *gin.Context)
}

type templeteMongoController struct {
	templeteService templete_mongo.Service
}

func NewTempleteMongoController(
	templeteService templete_mongo.Service,
) TempleteMongoController {
	return &templeteMongoController{
		templeteService: templeteService,
	}
}

func (c *templeteMongoController) GetAllTempletes(ginContext *gin.Context) {
	var requestModel requestmodels.GetAllTempleteRequest

	if err := ginContext.ShouldBindQuery(&requestModel); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := middlewares.GetAccount(ginContext)

	switch result, failure := c.templeteService.GetAllTempletes(templete_mongo.GetAllTempletesModel{
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
		viewModel := []contracts.TempleteMongoContract{}
		for _, t := range result.Templetes {
			viewModel = append(viewModel, contracts.CreateTempleteMongoContract(t))
		}
		ginContext.JSON(http.StatusOK, viewModel)
	}
}

func (c *templeteMongoController) GetTempleteCount(ginContext *gin.Context) {
	var requestModel requestmodels.GetTempleteCountRequest

	if err := ginContext.ShouldBindQuery(&requestModel); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := middlewares.GetAccount(ginContext)

	switch result, failure := c.templeteService.GetTempleteCount(templete_mongo.GetTempleteCountModel{
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

func (c *templeteMongoController) GetTempleteByID(ginContext *gin.Context) {
	var id string
	if id = ginContext.Param("id"); id == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	account := middlewares.GetAccount(ginContext)

	switch result, failure := c.templeteService.GetTempleteByID(templete_mongo.GetTempleteByIDModel{
		ID:             id,
		OrganizationID: account.ActiveOrganization.ID,
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	case result.TempleteNotFound:
		ginContext.JSON(http.StatusConflict, gin.H{"templete_not_found": result.TempleteNotFound})
	default:
		ginContext.JSON(http.StatusOK, contracts.CreateTempleteMongoContract(result.Templete))
	}
}

func (c *templeteMongoController) CreateTemplete(ginContext *gin.Context) {
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

	switch result, failure := c.templeteService.CreateTemplete(templete_mongo.CreateTempleteModel{
		OrganizationID: account.ActiveOrganization.ID,
		Name:           requestModel.Name,
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	default:
		ginContext.JSON(http.StatusOK, contracts.CreateTempleteMongoContract(result.Templete))
	}
}

func (c *templeteMongoController) UpdateTemplete(ginContext *gin.Context) {
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

	switch result, failure := c.templeteService.UpdateTemplete(templete_mongo.UpdateTempleteModel{
		ID:             id,
		OrganizationID: account.ActiveOrganization.ID,
		Name:           requestModel.Name,
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	case result.TempleteNotFound:
		ginContext.JSON(http.StatusConflict, gin.H{"templete_not_found": result.TempleteNotFound})
	default:
		ginContext.JSON(http.StatusOK, contracts.CreateTempleteMongoContract(result.Templete))
	}
}

func (c *templeteMongoController) DeleteTemplete(ginContext *gin.Context) {
	var id string
	if id = ginContext.Param("id"); id == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	account := middlewares.GetAccount(ginContext)

	switch result, failure := c.templeteService.DeleteTemplete(templete_mongo.DeleteTempleteModel{
		ID:             id,
		OrganizationID: account.ActiveOrganization.ID,
	}); {
	case failure != nil:
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": failure.Error()})
	case result.TempleteNotFound:
		ginContext.JSON(http.StatusConflict, gin.H{"templete_not_found": result.TempleteNotFound})
	default:
		ginContext.JSON(http.StatusOK, contracts.CreateTempleteMongoContract(result.Templete))
	}
}
