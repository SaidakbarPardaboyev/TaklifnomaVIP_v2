package swagger

// CreateOrderSwagger godoc
// @Summary      Create a new order
// @Tags         Order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  requestmodels.CreateOrderRequest  true  "Order creation payload"
// @Success      201  {object} contracts.OrderContract
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/order/create [post]
func CreateOrderSwagger() {}

// GetAllOrdersSwagger godoc
// @Summary      List orders for current account
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Param        status         query  string  false  "Filter by status"
// @Param        template_code  query  string  false  "Filter by template code"
// @Param        from_date      query  string  false  "Filter from date (RFC3339)"
// @Param        to_date        query  string  false  "Filter to date (RFC3339)"
// @Param        page           query  int     false  "Page number"
// @Param        limit          query  int     false  "Page size"
// @Param        sort_by        query  string  false  "Sort field"
// @Param        order          query  string  false  "Sort direction (asc/desc)"
// @Success      200  {array}  contracts.OrderContract
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/order/list [get]
func GetAllOrdersSwagger() {}

// GetOrderByIDSwagger godoc
// @Summary      Get order by ID
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "Order ID"
// @Success      200  {object} contracts.OrderContract
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/order/get/{id} [get]
func GetOrderByIDSwagger() {}

// UpdateOrderSwagger godoc
// @Summary      Update a draft order
// @Tags         Order
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                        true  "Order ID"
// @Param        body  body  requestmodels.UpdateOrderRequest  true  "Order update payload"
// @Success      200  {object} contracts.OrderContract
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      409  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/order/update/{id} [put]
func UpdateOrderSwagger() {}

// DeleteOrderSwagger godoc
// @Summary      Cancel/delete a draft order
// @Tags         Order
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "Order ID"
// @Success      204
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      409  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/order/delete/{id} [delete]
func DeleteOrderSwagger() {}
