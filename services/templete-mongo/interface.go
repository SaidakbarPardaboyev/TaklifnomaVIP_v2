package templete_mongo

type Service interface {
	CreateTemplete(model CreateTempleteModel) (result *CreateTempleteResult, err error)
	UpdateTemplete(model UpdateTempleteModel) (result *UpdateTempleteResult, err error)
	DeleteTemplete(model DeleteTempleteModel) (result *DeleteTempleteResult, err error)
	GetAllTempletes(model GetAllTempletesModel) (result *GetAllTempletesResult, err error)
	GetTempleteByID(model GetTempleteByIDModel) (result *GetTempleteByIDResult, err error)
	GetTempleteCount(model GetTempleteCountModel) (result *GetTempleteCountResult, err error)
}
