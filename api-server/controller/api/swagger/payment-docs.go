package swagger

// GeneratePaymeLinkSwagger godoc
// @Summary      Generate a Payme payment link for an order
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  requestmodels.GeneratePaymeLinkRequest  true  "Order ID and amount"
// @Success      200  {object} contracts.PaymentLinkContract
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/payment/generate-link [post]
func GeneratePaymeLinkSwagger() {}

// PaymeWebhookSwagger godoc
// @Summary      Payme payment gateway webhook
// @Description  Handles Payme JSON-RPC calls: CheckPerformTransaction, CreateTransaction, PerformTransaction, CancelTransaction, CheckTransaction, GetStatement
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        body  body  object{method=string,params=object}  true  "Payme JSON-RPC request"
// @Success      200  {object} map[string]interface{}
// @Router       /payme [post]
func PaymeWebhookSwagger() {}
