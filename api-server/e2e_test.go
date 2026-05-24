package apiserver_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	api_ctrl "saidakbar.origin/api-server/controller/api"
	middlewares "saidakbar.origin/api-server/middleware"
	"saidakbar.origin/caching/models"
	"saidakbar.origin/db/mongo/entity"
	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/dto"
	account_service "saidakbar.origin/services/account"
	order_service "saidakbar.origin/services/order"
	payment_service "saidakbar.origin/services/payment"
	transaction_service "saidakbar.origin/services/transaction"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	testToken      = "test-token"
	paymeValidKey  = "payme-prod-key"
	testAccountID  = "550e8400-e29b-41d4-a716-446655440001"
	testOrderID    = "550e8400-e29b-41d4-a716-446655440002"
)

// ─── Mock: TokenCache ─────────────────────────────────────────────────────────

type mockTokenCache struct{}

func (m *mockTokenCache) GetAccount(token string) (*models.Account, error) {
	if token == testToken {
		return testAccount(), nil
	}
	return nil, errors.New("unauthorized")
}

func (m *mockTokenCache) GetClient(_ string) (*models.ClientToken, error) {
	return nil, errors.New("not found")
}

func (m *mockTokenCache) SetAccount(_ string, _ *models.Account, _ time.Duration) error {
	return nil
}

// ─── Mock: AccountService ────────────────────────────────────────────────────

type mockAccountService struct {
	VerifyCodeFn func(account_service.VerifyCodeModel) (*dto.VerifyCodeResult, error)
}

func (m *mockAccountService) VerifyCode(model account_service.VerifyCodeModel) (*dto.VerifyCodeResult, error) {
	if m.VerifyCodeFn == nil {
		return &dto.VerifyCodeResult{}, nil
	}
	return m.VerifyCodeFn(model)
}

// ─── Mock: OrderService ───────────────────────────────────────────────────────

type mockOrderService struct {
	CreateFn             func(order_service.CreateOrderModel) (*order_service.CreateOrderResult, error)
	GetByIDFn            func(order_service.GetOrderByIDModel) (*order_service.GetOrderByIDResult, error)
	GetPublicFn          func(string) (*order_service.GetOrderByIDResult, error)
	GetAllFn             func(order_service.GetAllOrdersModel) (*order_service.GetAllOrdersResult, error)
	UpdateFn             func(order_service.UpdateOrderModel) (*order_service.UpdateOrderResult, error)
	DeleteFn             func(order_service.DeleteOrderModel) (*order_service.DeleteOrderResult, error)
	IncrementViewCountFn func(string) error
}

func (m *mockOrderService) Create(model order_service.CreateOrderModel) (*order_service.CreateOrderResult, error) {
	if m.CreateFn == nil {
		return nil, nil
	}
	return m.CreateFn(model)
}
func (m *mockOrderService) GetByID(model order_service.GetOrderByIDModel) (*order_service.GetOrderByIDResult, error) {
	if m.GetByIDFn == nil {
		return nil, nil
	}
	return m.GetByIDFn(model)
}
func (m *mockOrderService) GetPublic(id string) (*order_service.GetOrderByIDResult, error) {
	if m.GetPublicFn == nil {
		return nil, nil
	}
	return m.GetPublicFn(id)
}
func (m *mockOrderService) GetAll(model order_service.GetAllOrdersModel) (*order_service.GetAllOrdersResult, error) {
	if m.GetAllFn == nil {
		return &order_service.GetAllOrdersResult{Orders: []*entity.OrderEntity{}}, nil
	}
	return m.GetAllFn(model)
}
func (m *mockOrderService) Update(model order_service.UpdateOrderModel) (*order_service.UpdateOrderResult, error) {
	if m.UpdateFn == nil {
		return nil, nil
	}
	return m.UpdateFn(model)
}
func (m *mockOrderService) Delete(model order_service.DeleteOrderModel) (*order_service.DeleteOrderResult, error) {
	if m.DeleteFn == nil {
		return nil, nil
	}
	return m.DeleteFn(model)
}
func (m *mockOrderService) IncrementViewCount(id string) error {
	if m.IncrementViewCountFn == nil {
		return nil
	}
	return m.IncrementViewCountFn(id)
}

// ─── Mock: PaymentService ─────────────────────────────────────────────────────

type mockPaymentService struct {
	GeneratePaymeLinkFn       func(payment_service.GeneratePaymeLinkModel) (*payment_service.GeneratePaymeLinkResult, error)
	CheckPerformTransactionFn func(payment_service.CheckPerformTransactionModel) (*payment_service.CheckPerformTransactionResult, error)
	CreateTransactionFn       func(payment_service.CreateTransactionModel) (*payment_service.CreateTransactionResult, error)
	PerformTransactionFn      func(payment_service.PerformTransactionModel) (*payment_service.PerformTransactionResult, error)
	CancelTransactionFn       func(payment_service.CancelTransactionModel) (*payment_service.CancelTransactionResult, error)
	CheckTransactionFn        func(payment_service.CheckTransactionModel) (*payment_service.CheckTransactionResult, error)
	GetStatementFn            func(payment_service.GetStatementModel) (*payment_service.GetStatementResult, error)
}

func (m *mockPaymentService) GeneratePaymeLink(model payment_service.GeneratePaymeLinkModel) (*payment_service.GeneratePaymeLinkResult, error) {
	if m.GeneratePaymeLinkFn == nil {
		return nil, nil
	}
	return m.GeneratePaymeLinkFn(model)
}
func (m *mockPaymentService) CheckPerformTransaction(model payment_service.CheckPerformTransactionModel) (*payment_service.CheckPerformTransactionResult, error) {
	if m.CheckPerformTransactionFn == nil {
		return nil, nil
	}
	return m.CheckPerformTransactionFn(model)
}
func (m *mockPaymentService) CreateTransaction(model payment_service.CreateTransactionModel) (*payment_service.CreateTransactionResult, error) {
	if m.CreateTransactionFn == nil {
		return nil, nil
	}
	return m.CreateTransactionFn(model)
}
func (m *mockPaymentService) PerformTransaction(model payment_service.PerformTransactionModel) (*payment_service.PerformTransactionResult, error) {
	if m.PerformTransactionFn == nil {
		return nil, nil
	}
	return m.PerformTransactionFn(model)
}
func (m *mockPaymentService) CancelTransaction(model payment_service.CancelTransactionModel) (*payment_service.CancelTransactionResult, error) {
	if m.CancelTransactionFn == nil {
		return nil, nil
	}
	return m.CancelTransactionFn(model)
}
func (m *mockPaymentService) CheckTransaction(model payment_service.CheckTransactionModel) (*payment_service.CheckTransactionResult, error) {
	if m.CheckTransactionFn == nil {
		return nil, nil
	}
	return m.CheckTransactionFn(model)
}
func (m *mockPaymentService) GetStatement(model payment_service.GetStatementModel) (*payment_service.GetStatementResult, error) {
	if m.GetStatementFn == nil {
		return &payment_service.GetStatementResult{Transactions: []*mysql_entity.TransactionModel{}}, nil
	}
	return m.GetStatementFn(model)
}

// ─── Mock: TransactionService ─────────────────────────────────────────────────

type mockTransactionService struct {
	GetByOrderIDFn func(transaction_service.GetByOrderIDModel) (*transaction_service.GetByOrderIDResult, error)
	GetListFn      func(transaction_service.GetListModel) (*transaction_service.GetListResult, error)
}

func (m *mockTransactionService) GetByOrderID(model transaction_service.GetByOrderIDModel) (*transaction_service.GetByOrderIDResult, error) {
	if m.GetByOrderIDFn == nil {
		return nil, nil
	}
	return m.GetByOrderIDFn(model)
}
func (m *mockTransactionService) GetList(model transaction_service.GetListModel) (*transaction_service.GetListResult, error) {
	if m.GetListFn == nil {
		return &transaction_service.GetListResult{Transactions: []*mysql_entity.TransactionModel{}}, nil
	}
	return m.GetListFn(model)
}

// ─── Sample data ──────────────────────────────────────────────────────────────

func testAccount() *models.Account {
	return &models.Account{
		ID:        primitive.NewObjectID(),
		Username:  "+998901234567",
		Name:      "Ali Valiyev",
		TokenType: 0,
		ActiveOrganization: &models.Organization{
			ID:   testAccountID,
			Name: "Ali Valiyev",
		},
	}
}

func sampleOrder() *entity.OrderEntity {
	now := time.Now()
	return &entity.OrderEntity{
		ID:           testOrderID,
		AccountID:    testAccountID,
		TemplateCode: "wedding-001",
		InvitationInfo: &entity.InvitationInfo{
			Title:         "Ali & Zulfiya",
			GroomFullname: "Ali Valiyev",
			BrideFullname: "Zulfiya Karimova",
			PartnersNames: "Ali & Zulfiya",
			EventDate:     now,
			EventTime:     "15:00",
			Address:       "Toshkent, Yunusobod",
			Location:      "https://maps.example.com",
			Images:        []string{},
		},
		Price:     99000,
		Status:    entity.OrderStatusDraft,
		ViewCount: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func sampleTx() *mysql_entity.TransactionModel {
	now := time.Now()
	return &mysql_entity.TransactionModel{
		ID:          "550e8400-e29b-41d4-a716-446655440003",
		PaymentID:   "payme-pay-1",
		Amount:      9900000,
		CreatedTime: &now,
		State:       1,
		AccountID:   testAccountID,
		OrderID:     testOrderID,
	}
}

// ─── Router builder ───────────────────────────────────────────────────────────

func buildEngine(
	accountSvc account_service.Service,
	orderSvc order_service.Service,
	paymentSvc payment_service.Service,
	txSvc transaction_service.Service,
) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	auth := middlewares.AuthAccount(&mockTokenCache{})
	accountCtrl := api_ctrl.NewAccountController(accountSvc)

	r.POST("/account/verify-code", accountCtrl.VerifyCode)

	if orderSvc != nil {
		r.GET("/i/:id", func(ctx *gin.Context) {
			id := ctx.Param("id")
			_ = orderSvc.IncrementViewCount(id)
			result, err := orderSvc.GetPublic(id)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if result == nil || result.NotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "invitation not found"})
				return
			}
			ctx.JSON(http.StatusOK, result.Order)
		})
	}

	apiGrp := r.Group("/api", auth)
	apiGrp.GET("/account/get-me", accountCtrl.GetMe)

	v1 := r.Group("/v1", auth)

	if orderSvc != nil {
		oc := api_ctrl.NewOrderController(orderSvc)
		v1.POST("/orders", oc.Create)
		v1.GET("/orders", oc.GetAll)
		v1.GET("/orders/:id", oc.GetByID)
		v1.PUT("/orders/:id", oc.Update)
		v1.DELETE("/orders/:id", oc.Delete)
	}

	if paymentSvc != nil {
		pc := api_ctrl.NewPaymentController(paymentSvc)
		r.POST("/payme", middlewares.PaymeAuth(paymeValidKey, "payme-staging-key"), pc.PaymeWebhook)
		v1.POST("/payment/generate-link", pc.GeneratePaymeLink)
	}

	if txSvc != nil {
		tc := api_ctrl.NewTransactionController(txSvc)
		v1.GET("/transactions", tc.GetList)
		v1.GET("/transactions/:order_id", tc.GetByOrderID)
	}

	return r
}

// ─── Request helpers ──────────────────────────────────────────────────────────

func do(engine *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func doAuth(engine *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", testToken)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func doPayme(engine *gin.Engine, body string, validKey bool) *httptest.ResponseRecorder {
	key := paymeValidKey
	if !validKey {
		key = "wrong-key"
	}
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte("Paycom:"+key))
	req := httptest.NewRequest(http.MethodPost, "/payme", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func bodyJSON(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	require.NoError(t, json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&m))
	return m
}

// ─── Account tests ────────────────────────────────────────────────────────────

func TestVerifyCode(t *testing.T) {
	const path = "/account/verify-code"
	const validBody = `{"phone":"+998901234567","code":"123456"}`

	t.Run("success", func(t *testing.T) {
		svc := &mockAccountService{VerifyCodeFn: func(_ account_service.VerifyCodeModel) (*dto.VerifyCodeResult, error) {
			return &dto.VerifyCodeResult{Token: "my-token"}, nil
		}}
		w := do(buildEngine(svc, nil, nil, nil), http.MethodPost, path, validBody)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "my-token", bodyJSON(t, w)["token"])
	})

	t.Run("bad JSON", func(t *testing.T) {
		w := do(buildEngine(&mockAccountService{}, nil, nil, nil), http.MethodPost, path, `not-json`)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		svc := &mockAccountService{VerifyCodeFn: func(_ account_service.VerifyCodeModel) (*dto.VerifyCodeResult, error) {
			return &dto.VerifyCodeResult{UserNotFound: true}, nil
		}}
		w := do(buildEngine(svc, nil, nil, nil), http.MethodPost, path, validBody)
		assert.Equal(t, 404, w.Code)
	})

	t.Run("invalid code", func(t *testing.T) {
		svc := &mockAccountService{VerifyCodeFn: func(_ account_service.VerifyCodeModel) (*dto.VerifyCodeResult, error) {
			return &dto.VerifyCodeResult{InvalidCode: true}, nil
		}}
		w := do(buildEngine(svc, nil, nil, nil), http.MethodPost, path, validBody)
		assert.Equal(t, 400, w.Code)
		assert.Contains(t, bodyJSON(t, w)["error"], "invalid")
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAccountService{VerifyCodeFn: func(_ account_service.VerifyCodeModel) (*dto.VerifyCodeResult, error) {
			return nil, errors.New("db error")
		}}
		w := do(buildEngine(svc, nil, nil, nil), http.MethodPost, path, validBody)
		assert.Equal(t, 500, w.Code)
	})
}

func TestGetMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		w := doAuth(buildEngine(&mockAccountService{}, nil, nil, nil), http.MethodGet, "/api/account/get-me", "")
		assert.Equal(t, 200, w.Code)
		b := bodyJSON(t, w)
		assert.Equal(t, testAccountID, b["id"])
		assert.Equal(t, "Ali Valiyev", b["name"])
		assert.Equal(t, "+998901234567", b["phone"])
	})

	t.Run("no token", func(t *testing.T) {
		w := do(buildEngine(&mockAccountService{}, nil, nil, nil), http.MethodGet, "/api/account/get-me", "")
		assert.Equal(t, 401, w.Code)
	})
}

// ─── Order tests ──────────────────────────────────────────────────────────────

const validCreateOrderBody = `{
  "template_code": "wedding-001",
  "invitation_info": {
    "title": "Ali & Zulfiya",
    "groom_fullname": "Ali",
    "bride_fullname": "Zulfiya",
    "partners_names": "Ali & Zulfiya",
    "event_date": "2026-08-15T00:00:00Z",
    "event_time": "15:00",
    "address": "Toshkent",
    "location": "https://maps.example.com"
  }
}`

const validUpdateOrderBody = `{
  "invitation_info": {
    "title": "Ali & Zulfiya Updated",
    "event_time": "16:00",
    "address": "Samarqand",
    "location": "https://maps.example.com/2"
  }
}`

func TestCreateOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockOrderService{CreateFn: func(_ order_service.CreateOrderModel) (*order_service.CreateOrderResult, error) {
			return &order_service.CreateOrderResult{Order: sampleOrder()}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodPost, "/v1/orders", validCreateOrderBody)
		assert.Equal(t, 201, w.Code)
		b := bodyJSON(t, w)
		assert.Equal(t, testOrderID, b["id"])
		assert.Equal(t, "wedding-001", b["template_code"])
		assert.Equal(t, "draft", b["status"])
	})

	t.Run("missing template_code", func(t *testing.T) {
		body := `{"invitation_info":{"title":"T","groom_fullname":"G","bride_fullname":"B","event_time":"15:00","address":"A","location":"L"}}`
		w := doAuth(buildEngine(&mockAccountService{}, &mockOrderService{}, nil, nil), http.MethodPost, "/v1/orders", body)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockOrderService{CreateFn: func(_ order_service.CreateOrderModel) (*order_service.CreateOrderResult, error) {
			return nil, errors.New("db error")
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodPost, "/v1/orders", validCreateOrderBody)
		assert.Equal(t, 500, w.Code)
	})
}

func TestGetAllOrders(t *testing.T) {
	t.Run("success no filters", func(t *testing.T) {
		svc := &mockOrderService{GetAllFn: func(_ order_service.GetAllOrdersModel) (*order_service.GetAllOrdersResult, error) {
			return &order_service.GetAllOrdersResult{Orders: []*entity.OrderEntity{sampleOrder()}}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodGet, "/v1/orders", "")
		assert.Equal(t, 200, w.Code)
		var orders []interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &orders))
		assert.Len(t, orders, 1)
	})

	t.Run("with status filter", func(t *testing.T) {
		var capturedModel order_service.GetAllOrdersModel
		svc := &mockOrderService{GetAllFn: func(m order_service.GetAllOrdersModel) (*order_service.GetAllOrdersResult, error) {
			capturedModel = m
			return &order_service.GetAllOrdersResult{Orders: []*entity.OrderEntity{}}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodGet, "/v1/orders?status=active&page=2&limit=5", "")
		assert.Equal(t, 200, w.Code)
		require.NotNil(t, capturedModel.Status)
		assert.Equal(t, "active", *capturedModel.Status)
	})

	t.Run("invalid from_date format", func(t *testing.T) {
		w := doAuth(buildEngine(&mockAccountService{}, &mockOrderService{}, nil, nil), http.MethodGet, "/v1/orders?from_date=not-a-date", "")
		assert.Equal(t, 400, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockOrderService{GetAllFn: func(_ order_service.GetAllOrdersModel) (*order_service.GetAllOrdersResult, error) {
			return nil, errors.New("db error")
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodGet, "/v1/orders", "")
		assert.Equal(t, 500, w.Code)
	})
}

func TestGetOrderByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockOrderService{GetByIDFn: func(_ order_service.GetOrderByIDModel) (*order_service.GetOrderByIDResult, error) {
			return &order_service.GetOrderByIDResult{Order: sampleOrder()}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodGet, "/v1/orders/"+testOrderID, "")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, testOrderID, bodyJSON(t, w)["id"])
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockOrderService{GetByIDFn: func(_ order_service.GetOrderByIDModel) (*order_service.GetOrderByIDResult, error) {
			return &order_service.GetOrderByIDResult{NotFound: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodGet, "/v1/orders/"+testOrderID, "")
		assert.Equal(t, 404, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &mockOrderService{GetByIDFn: func(_ order_service.GetOrderByIDModel) (*order_service.GetOrderByIDResult, error) {
			return &order_service.GetOrderByIDResult{Unauthorized: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodGet, "/v1/orders/"+testOrderID, "")
		assert.Equal(t, 403, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockOrderService{GetByIDFn: func(_ order_service.GetOrderByIDModel) (*order_service.GetOrderByIDResult, error) {
			return nil, errors.New("db error")
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodGet, "/v1/orders/"+testOrderID, "")
		assert.Equal(t, 500, w.Code)
	})
}

func TestUpdateOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockOrderService{UpdateFn: func(_ order_service.UpdateOrderModel) (*order_service.UpdateOrderResult, error) {
			return &order_service.UpdateOrderResult{Order: sampleOrder()}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodPut, "/v1/orders/"+testOrderID, validUpdateOrderBody)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, testOrderID, bodyJSON(t, w)["id"])
	})

	t.Run("missing title", func(t *testing.T) {
		body := `{"invitation_info":{"event_time":"15:00","address":"A","location":"L"}}`
		w := doAuth(buildEngine(&mockAccountService{}, &mockOrderService{}, nil, nil), http.MethodPut, "/v1/orders/"+testOrderID, body)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockOrderService{UpdateFn: func(_ order_service.UpdateOrderModel) (*order_service.UpdateOrderResult, error) {
			return &order_service.UpdateOrderResult{NotFound: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodPut, "/v1/orders/"+testOrderID, validUpdateOrderBody)
		assert.Equal(t, 404, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &mockOrderService{UpdateFn: func(_ order_service.UpdateOrderModel) (*order_service.UpdateOrderResult, error) {
			return &order_service.UpdateOrderResult{Unauthorized: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodPut, "/v1/orders/"+testOrderID, validUpdateOrderBody)
		assert.Equal(t, 403, w.Code)
	})

	t.Run("not draft", func(t *testing.T) {
		svc := &mockOrderService{UpdateFn: func(_ order_service.UpdateOrderModel) (*order_service.UpdateOrderResult, error) {
			return &order_service.UpdateOrderResult{NotDraft: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodPut, "/v1/orders/"+testOrderID, validUpdateOrderBody)
		assert.Equal(t, 409, w.Code)
	})
}

func TestDeleteOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockOrderService{DeleteFn: func(_ order_service.DeleteOrderModel) (*order_service.DeleteOrderResult, error) {
			return &order_service.DeleteOrderResult{}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodDelete, "/v1/orders/"+testOrderID, "")
		assert.Equal(t, 204, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockOrderService{DeleteFn: func(_ order_service.DeleteOrderModel) (*order_service.DeleteOrderResult, error) {
			return &order_service.DeleteOrderResult{NotFound: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodDelete, "/v1/orders/"+testOrderID, "")
		assert.Equal(t, 404, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &mockOrderService{DeleteFn: func(_ order_service.DeleteOrderModel) (*order_service.DeleteOrderResult, error) {
			return &order_service.DeleteOrderResult{Unauthorized: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodDelete, "/v1/orders/"+testOrderID, "")
		assert.Equal(t, 403, w.Code)
	})

	t.Run("not draft", func(t *testing.T) {
		svc := &mockOrderService{DeleteFn: func(_ order_service.DeleteOrderModel) (*order_service.DeleteOrderResult, error) {
			return &order_service.DeleteOrderResult{NotDraft: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodDelete, "/v1/orders/"+testOrderID, "")
		assert.Equal(t, 409, w.Code)
	})
}

func TestGetPublicOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var viewCountCalled bool
		svc := &mockOrderService{
			IncrementViewCountFn: func(_ string) error {
				viewCountCalled = true
				return nil
			},
			GetPublicFn: func(_ string) (*order_service.GetOrderByIDResult, error) {
				return &order_service.GetOrderByIDResult{Order: sampleOrder()}, nil
			},
		}
		w := do(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodGet, "/i/order-uuid-1", "")
		assert.Equal(t, 200, w.Code)
		assert.True(t, viewCountCalled)
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockOrderService{
			GetPublicFn: func(_ string) (*order_service.GetOrderByIDResult, error) {
				return &order_service.GetOrderByIDResult{NotFound: true}, nil
			},
		}
		w := do(buildEngine(&mockAccountService{}, svc, nil, nil), http.MethodGet, "/i/order-uuid-1", "")
		assert.Equal(t, 404, w.Code)
	})
}

// ─── Payment tests ────────────────────────────────────────────────────────────

const validOrderUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestGeneratePaymeLink(t *testing.T) {
	validBody := `{"order_id":"` + validOrderUUID + `","amount":99000}`

	t.Run("success", func(t *testing.T) {
		link := "https://checkout.paycom.uz/abc"
		svc := &mockPaymentService{GeneratePaymeLinkFn: func(_ payment_service.GeneratePaymeLinkModel) (*payment_service.GeneratePaymeLinkResult, error) {
			return &payment_service.GeneratePaymeLinkResult{Link: &link}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, svc, nil), http.MethodPost, "/v1/payment/generate-link", validBody)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, link, bodyJSON(t, w)["link"])
	})

	t.Run("validation error — bad UUID", func(t *testing.T) {
		body := `{"order_id":"not-a-uuid","amount":99000}`
		w := doAuth(buildEngine(&mockAccountService{}, nil, &mockPaymentService{}, nil), http.MethodPost, "/v1/payment/generate-link", body)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("order not found", func(t *testing.T) {
		svc := &mockPaymentService{GeneratePaymeLinkFn: func(_ payment_service.GeneratePaymeLinkModel) (*payment_service.GeneratePaymeLinkResult, error) {
			return &payment_service.GeneratePaymeLinkResult{OrderNotFound: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, svc, nil), http.MethodPost, "/v1/payment/generate-link", validBody)
		assert.Equal(t, 404, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &mockPaymentService{GeneratePaymeLinkFn: func(_ payment_service.GeneratePaymeLinkModel) (*payment_service.GeneratePaymeLinkResult, error) {
			return &payment_service.GeneratePaymeLinkResult{Unauthorized: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, svc, nil), http.MethodPost, "/v1/payment/generate-link", validBody)
		assert.Equal(t, 403, w.Code)
	})

	t.Run("incorrect amount", func(t *testing.T) {
		svc := &mockPaymentService{GeneratePaymeLinkFn: func(_ payment_service.GeneratePaymeLinkModel) (*payment_service.GeneratePaymeLinkResult, error) {
			return &payment_service.GeneratePaymeLinkResult{IncorrectAmount: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, svc, nil), http.MethodPost, "/v1/payment/generate-link", validBody)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockPaymentService{GeneratePaymeLinkFn: func(_ payment_service.GeneratePaymeLinkModel) (*payment_service.GeneratePaymeLinkResult, error) {
			return nil, errors.New("db error")
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, svc, nil), http.MethodPost, "/v1/payment/generate-link", validBody)
		assert.Equal(t, 500, w.Code)
	})
}

func webhookBody(method, params string) string {
	return `{"method":"` + method + `","params":` + params + `}`
}

func TestPaymeWebhook(t *testing.T) {
	eng := func(svc *mockPaymentService) *gin.Engine {
		return buildEngine(&mockAccountService{}, nil, svc, nil)
	}

	t.Run("wrong auth key", func(t *testing.T) {
		w := doPayme(eng(&mockPaymentService{}), `{"method":"CheckPerformTransaction","params":{}}`, false)
		assert.Equal(t, 200, w.Code)
		b := bodyJSON(t, w)
		errObj := b["error"].(map[string]interface{})
		assert.Equal(t, float64(-32504), errObj["code"])
	})

	t.Run("bad JSON body", func(t *testing.T) {
		w := doPayme(eng(&mockPaymentService{}), `not-json`, true)
		assert.Equal(t, 200, w.Code)
		errObj := bodyJSON(t, w)["error"].(map[string]interface{})
		assert.Equal(t, float64(-31099), errObj["code"])
	})

	t.Run("unknown method", func(t *testing.T) {
		w := doPayme(eng(&mockPaymentService{}), `{"method":"UnknownMethod","params":{}}`, true)
		assert.Equal(t, 200, w.Code)
		errObj := bodyJSON(t, w)["error"].(map[string]interface{})
		assert.Equal(t, float64(-31099), errObj["code"])
	})

	// CheckPerformTransaction
	checkPerformParams := `{"account":{"account_id":"` + testAccountID + `","order_id":"` + testOrderID + `"},"amount":9900000}`

	t.Run("CheckPerformTransaction success", func(t *testing.T) {
		svc := &mockPaymentService{CheckPerformTransactionFn: func(_ payment_service.CheckPerformTransactionModel) (*payment_service.CheckPerformTransactionResult, error) {
			return &payment_service.CheckPerformTransactionResult{Succeed: true, Order: sampleOrder()}, nil
		}}
		w := doPayme(eng(svc), webhookBody("CheckPerformTransaction", checkPerformParams), true)
		assert.Equal(t, 200, w.Code)
		result := bodyJSON(t, w)["result"].(map[string]interface{})
		assert.Equal(t, true, result["allow"])
	})

	t.Run("CheckPerformTransaction not succeed", func(t *testing.T) {
		svc := &mockPaymentService{CheckPerformTransactionFn: func(_ payment_service.CheckPerformTransactionModel) (*payment_service.CheckPerformTransactionResult, error) {
			return &payment_service.CheckPerformTransactionResult{OrderNotFound: true, ErrorCode: -31050, MessageEn: "order not found"}, nil
		}}
		w := doPayme(eng(svc), webhookBody("CheckPerformTransaction", checkPerformParams), true)
		assert.Equal(t, 200, w.Code)
		errObj := bodyJSON(t, w)["error"].(map[string]interface{})
		assert.Equal(t, float64(-31050), errObj["code"])
	})

	t.Run("CheckPerformTransaction service error", func(t *testing.T) {
		svc := &mockPaymentService{CheckPerformTransactionFn: func(_ payment_service.CheckPerformTransactionModel) (*payment_service.CheckPerformTransactionResult, error) {
			return nil, errors.New("db error")
		}}
		w := doPayme(eng(svc), webhookBody("CheckPerformTransaction", checkPerformParams), true)
		assert.Equal(t, 200, w.Code)
		errObj := bodyJSON(t, w)["error"].(map[string]interface{})
		assert.Equal(t, float64(-32400), errObj["code"])
	})

	// CreateTransaction
	createTxParams := `{"account":{"account_id":"` + testAccountID + `","order_id":"` + testOrderID + `"},"id":"payme-pay-1","amount":9900000,"time":1716000000000}`

	t.Run("CreateTransaction success", func(t *testing.T) {
		svc := &mockPaymentService{CreateTransactionFn: func(_ payment_service.CreateTransactionModel) (*payment_service.CreateTransactionResult, error) {
			return &payment_service.CreateTransactionResult{Succeed: true, Transaction: sampleTx()}, nil
		}}
		w := doPayme(eng(svc), webhookBody("CreateTransaction", createTxParams), true)
		assert.Equal(t, 200, w.Code)
		result := bodyJSON(t, w)["result"].(map[string]interface{})
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440003", result["transaction"])
		assert.Equal(t, float64(1), result["state"])
	})

	t.Run("CreateTransaction not succeed", func(t *testing.T) {
		svc := &mockPaymentService{CreateTransactionFn: func(_ payment_service.CreateTransactionModel) (*payment_service.CreateTransactionResult, error) {
			return &payment_service.CreateTransactionResult{AlreadyPaid: true, ErrorCode: -31099, MessageEn: "already paid"}, nil
		}}
		w := doPayme(eng(svc), webhookBody("CreateTransaction", createTxParams), true)
		assert.Equal(t, 200, w.Code)
		errObj := bodyJSON(t, w)["error"].(map[string]interface{})
		assert.Equal(t, float64(-31099), errObj["code"])
	})

	// PerformTransaction
	performParams := `{"id":"payme-pay-1"}`

	t.Run("PerformTransaction success", func(t *testing.T) {
		tx := sampleTx()
		tx.State = 2
		now := time.Now()
		tx.PerformTime = &now
		svc := &mockPaymentService{PerformTransactionFn: func(_ payment_service.PerformTransactionModel) (*payment_service.PerformTransactionResult, error) {
			return &payment_service.PerformTransactionResult{Succeed: true, Transaction: tx}, nil
		}}
		w := doPayme(eng(svc), webhookBody("PerformTransaction", performParams), true)
		assert.Equal(t, 200, w.Code)
		result := bodyJSON(t, w)["result"].(map[string]interface{})
		assert.Equal(t, float64(2), result["state"])
		assert.NotEqual(t, float64(0), result["perform_time"])
	})

	t.Run("PerformTransaction not succeed", func(t *testing.T) {
		svc := &mockPaymentService{PerformTransactionFn: func(_ payment_service.PerformTransactionModel) (*payment_service.PerformTransactionResult, error) {
			return &payment_service.PerformTransactionResult{TransactionNotFound: true, ErrorCode: -31050, MessageEn: "not found"}, nil
		}}
		w := doPayme(eng(svc), webhookBody("PerformTransaction", performParams), true)
		assert.Equal(t, 200, w.Code)
		errObj := bodyJSON(t, w)["error"].(map[string]interface{})
		assert.Equal(t, float64(-31050), errObj["code"])
	})

	// CancelTransaction
	cancelParams := `{"id":"payme-pay-1","reason":3}`

	t.Run("CancelTransaction success", func(t *testing.T) {
		tx := sampleTx()
		tx.State = -1
		now := time.Now()
		tx.CancelTime = &now
		reason := 3
		tx.Reason = &reason
		svc := &mockPaymentService{CancelTransactionFn: func(_ payment_service.CancelTransactionModel) (*payment_service.CancelTransactionResult, error) {
			return &payment_service.CancelTransactionResult{Succeed: true, Transaction: tx}, nil
		}}
		w := doPayme(eng(svc), webhookBody("CancelTransaction", cancelParams), true)
		assert.Equal(t, 200, w.Code)
		result := bodyJSON(t, w)["result"].(map[string]interface{})
		assert.Equal(t, float64(-1), result["state"])
		assert.NotEqual(t, float64(0), result["cancel_time"])
	})

	t.Run("CancelTransaction not succeed", func(t *testing.T) {
		svc := &mockPaymentService{CancelTransactionFn: func(_ payment_service.CancelTransactionModel) (*payment_service.CancelTransactionResult, error) {
			return &payment_service.CancelTransactionResult{TransactionNotFound: true, ErrorCode: -31050, MessageEn: "not found"}, nil
		}}
		w := doPayme(eng(svc), webhookBody("CancelTransaction", cancelParams), true)
		assert.Equal(t, 200, w.Code)
		errObj := bodyJSON(t, w)["error"].(map[string]interface{})
		assert.Equal(t, float64(-31050), errObj["code"])
	})

	// CheckTransaction
	checkTxParams := `{"id":"payme-pay-1"}`

	t.Run("CheckTransaction success", func(t *testing.T) {
		svc := &mockPaymentService{CheckTransactionFn: func(_ payment_service.CheckTransactionModel) (*payment_service.CheckTransactionResult, error) {
			return &payment_service.CheckTransactionResult{Succeed: true, Transaction: sampleTx()}, nil
		}}
		w := doPayme(eng(svc), webhookBody("CheckTransaction", checkTxParams), true)
		assert.Equal(t, 200, w.Code)
		result := bodyJSON(t, w)["result"].(map[string]interface{})
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440003", result["transaction"])
		assert.Equal(t, float64(1), result["state"])
	})

	t.Run("CheckTransaction not succeed", func(t *testing.T) {
		svc := &mockPaymentService{CheckTransactionFn: func(_ payment_service.CheckTransactionModel) (*payment_service.CheckTransactionResult, error) {
			return &payment_service.CheckTransactionResult{TransactionNotFound: true, ErrorCode: -31050, MessageEn: "not found"}, nil
		}}
		w := doPayme(eng(svc), webhookBody("CheckTransaction", checkTxParams), true)
		assert.Equal(t, 200, w.Code)
		errObj := bodyJSON(t, w)["error"].(map[string]interface{})
		assert.Equal(t, float64(-31050), errObj["code"])
	})

	// GetStatement
	statementParams := `{"from":1716000000000,"to":1716100000000}`

	t.Run("GetStatement success", func(t *testing.T) {
		svc := &mockPaymentService{GetStatementFn: func(_ payment_service.GetStatementModel) (*payment_service.GetStatementResult, error) {
			return &payment_service.GetStatementResult{Transactions: []*mysql_entity.TransactionModel{sampleTx()}}, nil
		}}
		w := doPayme(eng(svc), webhookBody("GetStatement", statementParams), true)
		assert.Equal(t, 200, w.Code)
		result := bodyJSON(t, w)["result"].(map[string]interface{})
		txs := result["transactions"].([]interface{})
		assert.Len(t, txs, 1)
	})

	t.Run("GetStatement service error", func(t *testing.T) {
		svc := &mockPaymentService{GetStatementFn: func(_ payment_service.GetStatementModel) (*payment_service.GetStatementResult, error) {
			return nil, errors.New("db error")
		}}
		w := doPayme(eng(svc), webhookBody("GetStatement", statementParams), true)
		assert.Equal(t, 200, w.Code)
		errObj := bodyJSON(t, w)["error"].(map[string]interface{})
		assert.Equal(t, float64(-32400), errObj["code"])
	})
}

// ─── Transaction tests ────────────────────────────────────────────────────────

func TestGetTransactionList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockTransactionService{GetListFn: func(_ transaction_service.GetListModel) (*transaction_service.GetListResult, error) {
			return &transaction_service.GetListResult{Transactions: []*mysql_entity.TransactionModel{sampleTx()}}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, nil, svc), http.MethodGet, "/v1/transactions", "")
		assert.Equal(t, 200, w.Code)
		var txs []interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &txs))
		assert.Len(t, txs, 1)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockTransactionService{GetListFn: func(_ transaction_service.GetListModel) (*transaction_service.GetListResult, error) {
			return nil, errors.New("db error")
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, nil, svc), http.MethodGet, "/v1/transactions", "")
		assert.Equal(t, 500, w.Code)
	})
}

func TestGetTransactionByOrderID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockTransactionService{GetByOrderIDFn: func(_ transaction_service.GetByOrderIDModel) (*transaction_service.GetByOrderIDResult, error) {
			return &transaction_service.GetByOrderIDResult{Transaction: sampleTx()}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, nil, svc), http.MethodGet, "/v1/transactions/"+testOrderID, "")
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440003", bodyJSON(t, w)["id"])
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockTransactionService{GetByOrderIDFn: func(_ transaction_service.GetByOrderIDModel) (*transaction_service.GetByOrderIDResult, error) {
			return &transaction_service.GetByOrderIDResult{NotFound: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, nil, svc), http.MethodGet, "/v1/transactions/"+testOrderID, "")
		assert.Equal(t, 404, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		svc := &mockTransactionService{GetByOrderIDFn: func(_ transaction_service.GetByOrderIDModel) (*transaction_service.GetByOrderIDResult, error) {
			return &transaction_service.GetByOrderIDResult{Unauthorized: true}, nil
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, nil, svc), http.MethodGet, "/v1/transactions/"+testOrderID, "")
		assert.Equal(t, 403, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockTransactionService{GetByOrderIDFn: func(_ transaction_service.GetByOrderIDModel) (*transaction_service.GetByOrderIDResult, error) {
			return nil, errors.New("db error")
		}}
		w := doAuth(buildEngine(&mockAccountService{}, nil, nil, svc), http.MethodGet, "/v1/transactions/"+testOrderID, "")
		assert.Equal(t, 500, w.Code)
	})
}
