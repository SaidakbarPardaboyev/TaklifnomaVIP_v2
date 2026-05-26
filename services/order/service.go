package order_service

import (
	"time"

	"github.com/google/uuid"
	"saidakbar.origin/db/mongo/entity"
	"saidakbar.origin/plugins"
	"saidakbar.origin/repository"
	"saidakbar.origin/repository/models"
)

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) Service {
	return &orderService{orderRepo: orderRepo}
}

func (s *orderService) Create(model CreateOrderModel) (*CreateOrderResult, error) {
	now := plugins.GetNow()
	order := &entity.OrderEntity{
		ID:           uuid.NewString(),
		AccountID:    model.AccountID,
		TemplateCode: model.TemplateCode,
		InvitationInfo: &entity.InvitationInfo{
			Title:         model.InvitationInfo.Title,
			GroomFullname: model.InvitationInfo.GroomFullname,
			BrideFullname: model.InvitationInfo.BrideFullname,
			PartnersNames: model.InvitationInfo.PartnersNames,
			Description:   model.InvitationInfo.Description,
			EventDate:     model.InvitationInfo.EventDate,
			EventTime:     model.InvitationInfo.EventTime,
			Address:       model.InvitationInfo.Address,
			Location:      model.InvitationInfo.Location,
			VenueName:     model.InvitationInfo.VenueName,
			MainImage:     model.InvitationInfo.MainImage,
			Images:        model.InvitationInfo.Images,
			Story:         model.InvitationInfo.Story,
			Music:         model.InvitationInfo.Music,
		},
		Price:     model.Price,
		Status:    entity.OrderStatusDraft,
		ViewCount: 0,
		ExpiresAt: nil,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}
	return &CreateOrderResult{Order: order}, nil
}

func (s *orderService) GetByID(model GetOrderByIDModel) (*GetOrderByIDResult, error) {
	result := new(GetOrderByIDResult)

	order, err := s.orderRepo.GetByID(model.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		result.NotFound = true
		return result, nil
	}
	if order.AccountID != model.AccountID {
		result.Unauthorized = true
		return result, nil
	}

	result.Order = order
	return result, nil
}

func (s *orderService) GetPublic(id string) (*GetOrderByIDResult, error) {
	result := new(GetOrderByIDResult)
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if order == nil || order.Status != entity.OrderStatusActive {
		result.NotFound = true
		return result, nil
	}
	result.Order = order
	return result, nil
}

func (s *orderService) GetAll(model GetAllOrdersModel) (*GetAllOrdersResult, error) {
	result := new(GetAllOrdersResult)

	limit := 10
	if model.Limit != nil && *model.Limit > 0 && *model.Limit <= 100 {
		limit = *model.Limit
	}
	model.Limit = &limit

	var err error
	result.Orders, err = s.orderRepo.GetAll(&models.GetAllOrdersFilter{
		AccountID:    model.AccountID,
		Status:       model.Status,
		TemplateCode: model.TemplateCode,
		FromDate:     model.FromDate,
		ToDate:       model.ToDate,
		Page:         model.Page,
		Limit:        model.Limit,
		SortBy:       model.SortBy,
		Order:        model.Order,
	})
	return result, err
}

func (s *orderService) Update(model UpdateOrderModel) (*UpdateOrderResult, error) {
	result := new(UpdateOrderResult)

	order, err := s.orderRepo.GetByID(model.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		result.NotFound = true
		return result, nil
	}
	if order.AccountID != model.AccountID {
		result.Unauthorized = true
		return result, nil
	}
	if order.Status != entity.OrderStatusDraft {
		result.NotDraft = true
		return result, nil
	}

	order.InvitationInfo = &entity.InvitationInfo{
		Title:         model.InvitationInfo.Title,
		GroomFullname: model.InvitationInfo.GroomFullname,
		BrideFullname: model.InvitationInfo.BrideFullname,
		PartnersNames: model.InvitationInfo.PartnersNames,
		Description:   model.InvitationInfo.Description,
		EventDate:     model.InvitationInfo.EventDate,
		EventTime:     model.InvitationInfo.EventTime,
		Address:       model.InvitationInfo.Address,
		Location:      model.InvitationInfo.Location,
		VenueName:     model.InvitationInfo.VenueName,
		MainImage:     model.InvitationInfo.MainImage,
		Images:        model.InvitationInfo.Images,
		Story:         model.InvitationInfo.Story,
		Music:         model.InvitationInfo.Music,
	}
	order.UpdatedAt = time.Now()

	if err = s.orderRepo.Update(order); err != nil {
		return nil, err
	}
	result.Order = order
	return result, nil
}

func (s *orderService) Delete(model DeleteOrderModel) (*DeleteOrderResult, error) {
	result := new(DeleteOrderResult)

	order, err := s.orderRepo.GetByID(model.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		result.NotFound = true
		return result, nil
	}
	if order.AccountID != model.AccountID {
		result.Unauthorized = true
		return result, nil
	}
	if order.Status != entity.OrderStatusDraft {
		result.NotDraft = true
		return result, nil
	}

	return result, s.orderRepo.UpdateStatus(model.ID, entity.OrderStatusCancelled)
}

func (s *orderService) IncrementViewCount(id string) error {
	return s.orderRepo.IncrementViewCount(id)
}
