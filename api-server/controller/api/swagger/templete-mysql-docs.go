package swagger

// GetAllTemplateMysqlSwagger godoc
// @Summary      Get all templete mysql records
// @Tags         Templete-MySQL
// @Produce      json
// @Security     BearerAuth
// @Param        query  query  requestmodels.GetAllTempleteRequest  false  "Filter"
// @Success      200    {array}  contracts.TemplateMysqlContract
// @Failure      500    {object} map[string]string
// @Router       /api/templete-mysql/get [get]
func GetAllTemplateMysqlSwagger() {}

// GetTemplateMysqlCountSwagger godoc
// @Summary      Count templete mysql records
// @Tags         Templete-MySQL
// @Produce      json
// @Security     BearerAuth
// @Param        query  query  requestmodels.GetTempleteCountRequest  false  "Filter"
// @Success      200    {integer} int
// @Failure      500    {object}  map[string]string
// @Router       /api/templete-mysql/count [get]
func GetTemplateMysqlCountSwagger() {}

// GetTemplateMysqlByIDSwagger godoc
// @Summary      Get templete mysql record by ID
// @Tags         Templete-MySQL
// @Produce      json
// @Security     BearerAuth
// @Param        id   path   string  true  "Record ID"
// @Success      200  {object} contracts.TemplateMysqlContract
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/templete-mysql/get/{id} [get]
func GetTemplateMysqlByIDSwagger() {}

// CreateTemplateMysqlSwagger godoc
// @Summary      Create templete mysql record
// @Tags         Templete-MySQL
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  requestmodels.CreateTempleteRequest  true  "Record data"
// @Success      200   {object} contracts.TemplateMysqlContract
// @Failure      400   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /api/templete-mysql/create [post]
func CreateTemplateMysqlSwagger() {}

// UpdateTemplateMysqlSwagger godoc
// @Summary      Update templete mysql record
// @Tags         Templete-MySQL
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path   string                             true  "Record ID"
// @Param        body  body   requestmodels.UpdateTempleteRequest  true  "Record data"
// @Success      200   {object} contracts.TemplateMysqlContract
// @Failure      400   {object} map[string]string
// @Failure      404   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /api/templete-mysql/update/{id} [post]
func UpdateTemplateMysqlSwagger() {}

// DeleteTemplateMysqlSwagger godoc
// @Summary      Delete templete mysql record
// @Tags         Templete-MySQL
// @Produce      json
// @Security     BearerAuth
// @Param        id   path   string  true  "Record ID"
// @Success      200  {object} contracts.TemplateMysqlContract
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/templete-mysql/delete/{id} [post]
func DeleteTemplateMysqlSwagger() {}
