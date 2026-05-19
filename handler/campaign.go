package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
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

func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", 400, "error", errorMessage)
		c.JSON(400, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", 500, "error", nil)
		c.JSON(500, response)
		return
	}
	response := helper.APIResponse("Campaign created successfully", 200, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(200, response)
}

func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign", 400, "error", nil)
		c.JSON(400, response)
		return
	}
	var inputData campaign.CreateCampaignInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", 400, "error", errorMessage)
		c.JSON(400, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(input, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", 500, "error", nil)
		c.JSON(500, response)
		return
	}
	response := helper.APIResponse("Campaign updated successfully", 200, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(200, response)
}

func (h *CampaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput
	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to upload campaign image", 400, "error", errorMessage)
		c.JSON(400, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file")
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", 400, "error", nil)
		c.JSON(400, response)
		return
	}
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", 400, "error", nil)
		c.JSON(400, response)
		return
	}
	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		response := helper.APIResponse("Failed to upload campaign image", 400, "error", nil)
		c.JSON(400, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image uploaded successfully", 200, "success", data)
	c.JSON(200, response)
}
