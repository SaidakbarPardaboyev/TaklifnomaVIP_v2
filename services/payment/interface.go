package payment_service

type Service interface {
	GeneratePaymeLink(model GeneratePaymeLinkModel) (*GeneratePaymeLinkResult, error)
	CheckPerformTransaction(model CheckPerformTransactionModel) (*CheckPerformTransactionResult, error)
	CreateTransaction(model CreateTransactionModel) (*CreateTransactionResult, error)
	PerformTransaction(model PerformTransactionModel) (*PerformTransactionResult, error)
	CancelTransaction(model CancelTransactionModel) (*CancelTransactionResult, error)
	CheckTransaction(model CheckTransactionModel) (*CheckTransactionResult, error)
	GetStatement(model GetStatementModel) (*GetStatementResult, error)
}
