package handler

import (
	"bwastratup/helper"
	"bwastratup/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
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
	// token, err := h.jwtService.GenerateToken(newUser.ID)
	// if err != nil {
	// 	response := helper.APIResponse("Failed to generate token", 500, "error", nil)
	// 	c.JSON(500, response)
	// 	return
	// }
	formatter := user.FormatUser(newUser, "token")
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
	formatter := user.FormatUser(loggedinUser, "token")
	response := helper.APIResponse("Login successful", 200, "success", formatter)
	c.JSON(200, response)
}
