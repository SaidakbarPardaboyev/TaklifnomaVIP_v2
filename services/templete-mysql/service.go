package templete_mysql

import (
	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/plugins"
	"saidakbar.origin/repository"
	"saidakbar.origin/repository/models"

	"github.com/google/uuid"
)

type templatMysqlService struct {
	templeteRepository repository.TemplateMysqlRepository
	// TODO: inject other repositories or API clients here
}

func NewTemplateMysqlService(
	templeteRepository repository.TemplateMysqlRepository,
) Service {
	return &templatMysqlService{
		templeteRepository: templeteRepository,
	}
}

func (s *templatMysqlService) CreateTemplete(model CreateTempleteModel) (result *CreateTempleteResult, err error) {
	result = new(CreateTempleteResult)

	templete := &mysql_entity.TempleteModel{
		BaseEntity: mysql_entity.BaseEntity{
			ID:        uuid.New().String(),
			CreatedAt: plugins.GetNow(),
			UpdatedAt: plugins.GetNow(),
			IsDeleted: false,
		},
		OrganizationID: model.OrganizationID,
		Name:           model.Name,
	}

	if err = s.templeteRepository.CreateTemplete(templete); err != nil {
		return
	}

	result.Templete = templete
	return
}

func (s *templatMysqlService) UpdateTemplete(model UpdateTempleteModel) (result *UpdateTempleteResult, err error) {
	result = new(UpdateTempleteResult)

	templete, err := s.templeteRepository.GetTemplete(model.ID, model.OrganizationID)
	if err != nil {
		return
	}
	if templete == nil {
		result.TempleteNotFound = true
		return
	}

	templete.Name = model.Name
	templete.UpdatedAt = plugins.GetNow()

	if err = s.templeteRepository.UpdateTemplete(templete); err != nil {
		return
	}

	result.Templete = templete
	return
}

func (s *templatMysqlService) DeleteTemplete(model DeleteTempleteModel) (result *DeleteTempleteResult, err error) {
	result = new(DeleteTempleteResult)

	templete, err := s.templeteRepository.GetTemplete(model.ID, model.OrganizationID)
	if err != nil {
		return
	}
	if templete == nil {
		result.TempleteNotFound = true
		return
	}

	now := plugins.GetNow()
	templete.IsDeleted = true
	templete.DeletedAt = &now

	if err = s.templeteRepository.DeleteTemplete(templete); err != nil {
		return
	}

	result.Templete = templete
	return
}

func (s *templatMysqlService) GetAllTempletes(model GetAllTempletesModel) (result *GetAllTempletesResult, err error) {
	result = new(GetAllTempletesResult)

	result.Templetes, err = s.templeteRepository.GetAllTempletes(&models.GetAllTempleteFilter{
		OrganizationID: model.OrganizationID,
		Skip:           model.Skip,
		Limit:          model.Limit,
		Search:         model.Search,
		StartDate:      model.StartDate,
		EndDate:        model.EndDate,
	})
	return
}

func (s *templatMysqlService) GetTempleteByID(model GetTempleteByIDModel) (result *GetTempleteByIDResult, err error) {
	result = new(GetTempleteByIDResult)

	templete, err := s.templeteRepository.GetTemplete(model.ID, model.OrganizationID)
	if err != nil {
		return
	}
	if templete == nil {
		result.TempleteNotFound = true
		return
	}

	result.Templete = templete
	return
}

func (s *templatMysqlService) GetTempleteCount(model GetTempleteCountModel) (result *GetTempleteCountResult, err error) {
	result = new(GetTempleteCountResult)

	result.Count, err = s.templeteRepository.GetTempleteCount(&models.GetTempleteCountFilter{
		OrganizationID: model.OrganizationID,
		Search:         model.Search,
		StartDate:      model.StartDate,
		EndDate:        model.EndDate,
	})
	return
}
