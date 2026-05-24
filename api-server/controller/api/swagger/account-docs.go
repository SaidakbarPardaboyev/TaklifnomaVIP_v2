package swagger

// VerifyCodeSwagger godoc
// @Summary      Verify OTP code
// @Description  Submit the phone number and OTP code received via Telegram bot to get a bearer token
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        body  body  requestmodels.VerifyCodeRequest  true  "Phone and OTP code"
// @Success      200   {object} map[string]string
// @Failure      400   {object} map[string]string
// @Failure      404   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /account/verify-code [post]
func VerifyCodeSwagger() {}

// GetMeSwagger godoc
// @Summary      Get current account
// @Tags         Account
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object} contracts.AccountContract
// @Failure      401  {object} map[string]string
// @Router       /api/account/get-me [get]
func GetMeSwagger() {}
