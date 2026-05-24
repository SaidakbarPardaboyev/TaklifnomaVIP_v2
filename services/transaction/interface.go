package transaction_service

type Service interface {
	GetByOrderID(model GetByOrderIDModel) (*GetByOrderIDResult, error)
	GetList(model GetListModel) (*GetListResult, error)
}
