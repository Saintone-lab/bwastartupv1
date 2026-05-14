package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
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

func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Error to get campaign", 400, "error", nil)
		c.JSON(400, response)
		return
	}
	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Error to get campaign", 400, "error", nil)
		c.JSON(400, response)
		return
	}
	response := helper.APIResponse("Campaign detail", 200, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(200, response)
}
