package contracts

import (
	"saidakbar.origin/core"
	"saidakbar.origin/db/mongo/entity"
	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/plugins"
	payment_service "saidakbar.origin/services/payment"
)

// region Payme error
type paymeMessageDetail struct {
	Uz string `json:"uz"`
	Ru string `json:"ru"`
	En string `json:"en"`
}

type paymeErrorDetail struct {
	Code    int                `json:"code"`
	Message paymeMessageDetail `json:"message"`
}

type PaymeErrorContract struct {
	Error paymeErrorDetail `json:"error"`
}

func CreatePaymeErrorContract(code int, uz, ru, en string) PaymeErrorContract {
	return PaymeErrorContract{Error: paymeErrorDetail{
		Code:    code,
		Message: paymeMessageDetail{Uz: uz, Ru: ru, En: en},
	}}
}

// region Payment link
type PaymentLinkContract struct {
	Link string `json:"link"`
}

func CreatePaymentLinkContract(result payment_service.GeneratePaymeLinkResult) PaymentLinkContract {
	return PaymentLinkContract{Link: *result.Link}
}

// region CheckPerformTransaction
type checkPerformResult struct {
	Allow  bool                   `json:"allow"`
	Detail *receiptDetail         `json:"detail,omitempty"`
	Additional map[string]interface{} `json:"additional,omitempty"`
}

type receiptDetail struct {
	ReceiptType int           `json:"receipt_type"`
	Items       []receiptItem `json:"items"`
}

type receiptItem struct {
	Discount    int64  `json:"discount"`
	Title       string `json:"title"`
	Price       int64  `json:"price"`
	Count       int    `json:"count"`
	Code        string `json:"code"`
	Units       int    `json:"units"`
	VATPercent  int    `json:"vat_percent"`
	PackageCode string `json:"package_code"`
}

type CheckPerformTransactionContract struct {
	Result *checkPerformResult `json:"result"`
}

func CreateCheckPerformTransactionContract(order *entity.OrderEntity) CheckPerformTransactionContract {
	return CheckPerformTransactionContract{
		Result: &checkPerformResult{
			Allow:      true,
			Additional: map[string]interface{}{},
			Detail: &receiptDetail{
				ReceiptType: core.PaymeReceiptTypeSale,
				Items: []receiptItem{
					{
						Title:       "Wedding invitation — " + order.TemplateCode,
						Price:       int64(plugins.ConvertToPaymeAmount(order.Price)),
						Count:       1,
						Code:        "10399001001000000",
						Discount:    0,
						Units:       0,
						VATPercent:  0,
						PackageCode: "1504176",
					},
				},
			},
		},
	}
}

// region CreateTransaction
type CreateTransactionContract struct {
	Result struct {
		CreateTime  int64  `json:"create_time"`
		Transaction string `json:"transaction"`
		State       int    `json:"state"`
	} `json:"result"`
}

func BuildCreateTransactionContract(t *mysql_entity.TransactionModel) CreateTransactionContract {
	return CreateTransactionContract{Result: struct {
		CreateTime  int64  `json:"create_time"`
		Transaction string `json:"transaction"`
		State       int    `json:"state"`
	}{
		CreateTime:  plugins.ConverterToPaymeTimeFormat(t.CreatedTime),
		Transaction: t.ID,
		State:       t.State,
	}}
}

// region PerformTransaction
type PerformTransactionContract struct {
	Result struct {
		PerformTime int64  `json:"perform_time"`
		Transaction string `json:"transaction"`
		State       int    `json:"state"`
	} `json:"result"`
}

func BuildPerformTransactionContract(t *mysql_entity.TransactionModel) PerformTransactionContract {
	return PerformTransactionContract{Result: struct {
		PerformTime int64  `json:"perform_time"`
		Transaction string `json:"transaction"`
		State       int    `json:"state"`
	}{
		PerformTime: plugins.ConverterToPaymeTimeFormat(t.PerformTime),
		Transaction: t.ID,
		State:       t.State,
	}}
}

// region CancelTransaction
type CancelTransactionContract struct {
	Result struct {
		CancelTime  int64  `json:"cancel_time"`
		Transaction string `json:"transaction"`
		State       int    `json:"state"`
	} `json:"result"`
}

func BuildCancelTransactionContract(t *mysql_entity.TransactionModel) CancelTransactionContract {
	return CancelTransactionContract{Result: struct {
		CancelTime  int64  `json:"cancel_time"`
		Transaction string `json:"transaction"`
		State       int    `json:"state"`
	}{
		CancelTime:  plugins.ConverterToPaymeTimeFormat(t.CancelTime),
		Transaction: t.ID,
		State:       t.State,
	}}
}

// region CheckTransaction
type CheckTransactionContract struct {
	Result struct {
		CreateTime  int64  `json:"create_time"`
		PerformTime int64  `json:"perform_time"`
		CancelTime  int64  `json:"cancel_time"`
		Transaction string `json:"transaction"`
		State       int    `json:"state"`
		Reason      *int   `json:"reason"`
	} `json:"result"`
}

func BuildCheckTransactionContract(t *mysql_entity.TransactionModel) CheckTransactionContract {
	return CheckTransactionContract{Result: struct {
		CreateTime  int64  `json:"create_time"`
		PerformTime int64  `json:"perform_time"`
		CancelTime  int64  `json:"cancel_time"`
		Transaction string `json:"transaction"`
		State       int    `json:"state"`
		Reason      *int   `json:"reason"`
	}{
		CreateTime:  plugins.ConverterToPaymeTimeFormat(t.CreatedTime),
		PerformTime: plugins.ConverterToPaymeTimeFormat(t.PerformTime),
		CancelTime:  plugins.ConverterToPaymeTimeFormat(t.CancelTime),
		Transaction: t.ID,
		State:       t.State,
		Reason:      t.Reason,
	}}
}

// region GetStatement
type getStatementTransaction struct {
	ID          string  `json:"id"`
	Time        int64   `json:"time"`
	Amount      int     `json:"amount"`
	Account     account `json:"account"`
	CreateTime  int64   `json:"create_time"`
	PerformTime int64   `json:"perform_time"`
	CancelTime  int64   `json:"cancel_time"`
	Transaction string  `json:"transaction"`
	State       int     `json:"state"`
	Reason      *int    `json:"reason"`
}

type account struct {
	AccountID string `json:"account_id"`
	OrderID   string `json:"order_id"`
}

type GetStatementContract struct {
	Result struct {
		Transactions []getStatementTransaction `json:"transactions"`
	} `json:"result"`
}

func BuildGetStatementContract(txs []*mysql_entity.TransactionModel) GetStatementContract {
	items := make([]getStatementTransaction, 0, len(txs))
	for _, t := range txs {
		items = append(items, getStatementTransaction{
			ID:          t.ID,
			Time:        plugins.ConverterToPaymeTimeFormat(t.CreatedTime),
			Amount:      int(t.Amount),
			Account:     account{AccountID: t.AccountID, OrderID: t.OrderID},
			CreateTime:  plugins.ConverterToPaymeTimeFormat(t.CreatedTime),
			PerformTime: plugins.ConverterToPaymeTimeFormat(t.PerformTime),
			CancelTime:  plugins.ConverterToPaymeTimeFormat(t.CancelTime),
			Transaction: t.PaymentID,
			State:       t.State,
			Reason:      t.Reason,
		})
	}
	c := GetStatementContract{}
	c.Result.Transactions = items
	return c
}
