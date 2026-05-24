package templete_mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"saidakbar.origin/db/mongo/entity"
	"saidakbar.origin/plugins"
	"saidakbar.origin/repository"
	"saidakbar.origin/repository/models"
)

type templeteMongoService struct {
	templeteRepository repository.TempleteMongoRepository
	// TODO: inject other repositories or API clients here
}

func NewTempleteMongoService(
	templeteRepository repository.TempleteMongoRepository,
) Service {
	return &templeteMongoService{
		templeteRepository: templeteRepository,
	}
}

func (s *templeteMongoService) CreateTemplete(model CreateTempleteModel) (result *CreateTempleteResult, err error) {
	result = new(CreateTempleteResult)

	templete := &entity.TempleteEntity{
		BaseEntity: entity.BaseEntity{
			ID:        primitive.NewObjectID(),
			IsDeleted: false,
			CreatedAt: plugins.GetNow(),
			UpdatedAt: plugins.GetNow(),
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

func (s *templeteMongoService) UpdateTemplete(model UpdateTempleteModel) (result *UpdateTempleteResult, err error) {
	result = new(UpdateTempleteResult)

	var templete *entity.TempleteEntity
	{
		switch foundTemplete, failure := s.templeteRepository.GetTemplete(bson.M{
			entity.FieldTempleteID:             plugins.ObjectIDOrNil(model.ID),
			entity.FieldTempleteOrganizationID: model.OrganizationID,
			entity.FieldTempleteIsDeleted:      false,
		}); {
		case failure != nil:
			err = failure
			return
		case foundTemplete == nil:
			result.TempleteNotFound = true
			return
		default:
			templete = foundTemplete
		}
	}

	templete.Name = model.Name
	templete.UpdatedAt = plugins.GetNow()

	if err = s.templeteRepository.UpdateTemplete(templete); err != nil {
		return
	}

	result.Templete = templete
	return
}

func (s *templeteMongoService) DeleteTemplete(model DeleteTempleteModel) (result *DeleteTempleteResult, err error) {
	result = new(DeleteTempleteResult)

	var templete *entity.TempleteEntity
	{
		switch foundTemplete, failure := s.templeteRepository.GetTemplete(bson.M{
			entity.FieldTempleteID:             plugins.ObjectIDOrNil(model.ID),
			entity.FieldTempleteOrganizationID: model.OrganizationID,
			entity.FieldTempleteIsDeleted:      false,
		}); {
		case failure != nil:
			err = failure
			return
		case foundTemplete == nil:
			result.TempleteNotFound = true
			return
		default:
			templete = foundTemplete
		}
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

func (s *templeteMongoService) GetAllTempletes(model GetAllTempletesModel) (result *GetAllTempletesResult, err error) {
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

func (s *templeteMongoService) GetTempleteByID(model GetTempleteByIDModel) (result *GetTempleteByIDResult, err error) {
	result = new(GetTempleteByIDResult)

	switch foundTemplete, failure := s.templeteRepository.GetTemplete(bson.M{
		entity.FieldTempleteID:             plugins.ObjectIDOrNil(model.ID),
		entity.FieldTempleteOrganizationID: model.OrganizationID,
		entity.FieldTempleteIsDeleted:      false,
	}); {
	case failure != nil:
		err = failure
		return
	case foundTemplete == nil:
		result.TempleteNotFound = true
		return
	default:
		result.Templete = foundTemplete
	}

	return
}

func (s *templeteMongoService) GetTempleteCount(model GetTempleteCountModel) (result *GetTempleteCountResult, err error) {
	result = new(GetTempleteCountResult)

	result.Count, err = s.templeteRepository.GetTempleteCount(&models.GetTempleteCountFilter{
		OrganizationID: model.OrganizationID,
		Search:         model.Search,
		StartDate:      model.StartDate,
		EndDate:        model.EndDate,
	})
	return
}
