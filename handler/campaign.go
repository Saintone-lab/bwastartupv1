package handler

import (
	"bwastratup/campaign"
	"bwastratup/helper"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *CampaignHandler {
	return &CampaignHandler{service}
}

func (h *CampaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", 400, "error", nil)
		c.JSON(400, response)
		return
	}
	response := helper.APIResponse("List of campaigns", 200, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(200, response)
}
