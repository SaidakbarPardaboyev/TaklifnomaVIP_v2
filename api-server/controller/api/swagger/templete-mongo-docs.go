package swagger

// GetAllTempleteMongoSwagger godoc
// @Summary      Get all templete mongo records
// @Tags         Templete-Mongo
// @Produce      json
// @Security     BearerAuth
// @Param        query  query  requestmodels.GetAllTempleteRequest  false  "Filter"
// @Success      200    {array}  contracts.TempleteMongoContract
// @Failure      500    {object} map[string]string
// @Router       /api/templete-mongo/get [get]
func GetAllTempleteMongoSwagger() {}

// GetTempleteMongoCountSwagger godoc
// @Summary      Count templete mongo records
// @Tags         Templete-Mongo
// @Produce      json
// @Security     BearerAuth
// @Param        query  query  requestmodels.GetTempleteCountRequest  false  "Filter"
// @Success      200    {integer} int
// @Failure      500    {object}  map[string]string
// @Router       /api/templete-mongo/count [get]
func GetTempleteMongoCountSwagger() {}

// GetTempleteMongoByIDSwagger godoc
// @Summary      Get templete mongo record by ID
// @Tags         Templete-Mongo
// @Produce      json
// @Security     BearerAuth
// @Param        id   path   string  true  "Record ID"
// @Success      200  {object} contracts.TempleteMongoContract
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/templete-mongo/get/{id} [get]
func GetTempleteMongoByIDSwagger() {}

// CreateTempleteMongoSwagger godoc
// @Summary      Create templete mongo record
// @Tags         Templete-Mongo
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  requestmodels.CreateTempleteRequest  true  "Record data"
// @Success      200   {object} contracts.TempleteMongoContract
// @Failure      400   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /api/templete-mongo/create [post]
func CreateTempleteMongoSwagger() {}

// UpdateTempleteMongoSwagger godoc
// @Summary      Update templete mongo record
// @Tags         Templete-Mongo
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path   string                             true  "Record ID"
// @Param        body  body   requestmodels.UpdateTempleteRequest  true  "Record data"
// @Success      200   {object} contracts.TempleteMongoContract
// @Failure      400   {object} map[string]string
// @Failure      404   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /api/templete-mongo/update/{id} [post]
func UpdateTempleteMongoSwagger() {}

// DeleteTempleteMongoSwagger godoc
// @Summary      Delete templete mongo record
// @Tags         Templete-Mongo
// @Produce      json
// @Security     BearerAuth
// @Param        id   path   string  true  "Record ID"
// @Success      200  {object} contracts.TempleteMongoContract
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/templete-mongo/delete/{id} [post]
func DeleteTempleteMongoSwagger() {}
