package order_service

type Service interface {
	Create(model CreateOrderModel) (*CreateOrderResult, error)
	GetByID(model GetOrderByIDModel) (*GetOrderByIDResult, error)
	GetPublic(id string) (*GetOrderByIDResult, error)
	GetAll(model GetAllOrdersModel) (*GetAllOrdersResult, error)
	Update(model UpdateOrderModel) (*UpdateOrderResult, error)
	Delete(model DeleteOrderModel) (*DeleteOrderResult, error)
	IncrementViewCount(id string) error
}
