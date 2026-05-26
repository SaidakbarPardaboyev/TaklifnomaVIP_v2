package swagger

// GetTransactionListSwagger godoc
// @Summary      List transactions for current account
// @Tags         Transaction
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  contracts.TransactionContract
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/transaction/list [get]
func GetTransactionListSwagger() {}

// GetTransactionByOrderIDSwagger godoc
// @Summary      Get transaction by order ID
// @Tags         Transaction
// @Produce      json
// @Security     BearerAuth
// @Param        order_id  path  string  true  "Order ID"
// @Success      200  {object} contracts.TransactionContract
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/transaction/get/{order_id} [get]
func GetTransactionByOrderIDSwagger() {}
