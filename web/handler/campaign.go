package handler

import (
	"bwastartup/campaign"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService     user.Service
}

func NewCampaignHandler(campaignService campaign.Service, userService user.Service) *campaignHandler {
	return &campaignHandler{campaignService, userService}
}

func (h *campaignHandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

func (h *campaignHandler) New(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	input := campaign.FormCreateCampaignInput{}
	input.Users = users
	c.HTML(http.StatusOK, "campaign_new.html", input)
}

func (h *campaignHandler) Create(c *gin.Context) {
	var input campaign.FormCreateCampaignInput
	if err := c.ShouldBind(&input); err != nil {
		users, e := h.userService.GetAllUsers()
		if e != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": e.Error()})
			return
		}
		input.Users = users
		input.Error = err.Error()
		c.HTML(http.StatusOK, "campaign_new.html", input)
		return
	}

	user, err := h.userService.GetUserByID(input.UserID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	_, err = h.campaignService.CreateCampaign(campaign.CreateCampaignInput{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		User:             user,
	})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) NewImage(c *gin.Context) {
	campaignID := c.Param("id")
	id, _ := strconv.Atoi(campaignID)
	c.HTML(http.StatusOK, "campaign_image.html", gin.H{"ID": id})
}

func (h *campaignHandler) CreateImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	campaignID := c.Param("id")
	id, _ := strconv.Atoi(campaignID)

	existingCampaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	userID := existingCampaign.UserID

	fileName := file.Filename
	path := fmt.Sprintf("images/%d-%s", userID, fileName)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	userCampaign, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	_, err = h.campaignService.SaveCampaignImage(campaign.CreateCampaignImageInput{
		CampaignID: id,
		IsPrimary:  true,
		User:       userCampaign,
	}, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Edit(c *gin.Context) {
	campaignID := c.Param("id")
	id, _ := strconv.Atoi(campaignID)
	existingCampaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	input := campaign.FormUpdateCampaignInput{}
	input.ID = id
	input.Name = existingCampaign.Name
	input.ShortDescription = existingCampaign.ShortDescription
	input.Description = existingCampaign.Description
	input.GoalAmount = existingCampaign.GoalAmount
	input.Perks = existingCampaign.Perks
	c.HTML(http.StatusOK, "campaign_edit.html", input)
}

func (h *campaignHandler) Update(c *gin.Context) {
	campaignID := c.Param("id")
	id, _ := strconv.Atoi(campaignID)

	var input campaign.FormUpdateCampaignInput
	if err := c.ShouldBind(&input); err != nil {
		input.Error = err.Error()
		input.ID = id
		c.HTML(http.StatusOK, "campaign_edit.html", input)
		return
	}

	input.ID = id
	existingCampaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	userID := existingCampaign.UserID
	userCampaign, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	_, err = h.campaignService.UpdateCampaign(campaign.GetCampaignDetailInput{ID: id}, campaign.CreateCampaignInput{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		User:             userCampaign,
	})

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")
}

func (h *campaignHandler) Show(c *gin.Context) {
	campaignID := c.Param("id")
	id, _ := strconv.Atoi(campaignID)
	existingCampaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "campaign_show.html", existingCampaign)
}
