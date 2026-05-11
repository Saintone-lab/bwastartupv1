package handler

import (
	"bwastratup/auth"
	"bwastratup/helper"
	"bwastratup/user"
	"fmt"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {

		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}

		response := helper.APIResponse("Invalid input", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}
	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("Failed to register user", 500, "error", nil)
		c.JSON(500, response)
		return
	}
	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Failed to generate token", 500, "error", nil)
		c.JSON(500, response)
		return
	}
	formatter := user.FormatUser(newUser, token)
	response := helper.APIResponse("User registered successfully", 200, "success", formatter)
	c.JSON(200, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Invalid input", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}
	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": []string{err.Error()}}
		response := helper.APIResponse("Login failed", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}
	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.APIResponse("Failed to generate token", 500, "error", nil)
		c.JSON(500, response)
		return
	}
	formatter := user.FormatUser(loggedinUser, token)
	response := helper.APIResponse("Login successful", 200, "success", formatter)
	c.JSON(200, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Email Check Failed", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}
	isAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Failed to check email availability"}
		response := helper.APIResponse("Failed to check email availability", 500, "error", errorMessage)
		c.JSON(500, response)
		return
	}
	data := gin.H{
		"is_available": isAvailable,
	}

	metaMessage := "Email has been registered"
	if isAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, 200, "success", data)
	c.JSON(200, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar 1", 400, "error", data)
		c.JSON(400, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar 2", 400, "error", data)
		c.JSON(400, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar 3", 400, "error", data)
		c.JSON(400, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar uploaded successfully", 200, "success", data)
	c.JSON(200, response)
}
