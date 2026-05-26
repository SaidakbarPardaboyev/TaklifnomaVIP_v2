package swagger

// GetPublicInvitationSwagger godoc
// @Summary      Get public invitation by ID
// @Description  Returns the invitation details and increments view count. No authentication required.
// @Tags         Invitation
// @Produce      json
// @Param        id   path  string  true  "Invitation ID"
// @Success      200  {object} contracts.OrderContract
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /public/i/{id} [get]
func GetPublicInvitationSwagger() {}
