package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", 400, "error", nil)
		c.JSON(400, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", 500, "error", nil)
		c.JSON(500, response)
		return
	}
	response := helper.APIResponse("Campaign's transactions", 200, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(200, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	transactions, err := h.service.GetTransactionsByUserID(currentUser.ID)

	if err != nil {
		response := helper.APIResponse("Failed to get user's transactions", 500, "error", nil)
		c.JSON(500, response)
		return
	}

	response := helper.APIResponse("User's transactions", 200, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(200, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create transaction", 400, "error", errorMessage)
		c.JSON(400, response)
		return
	}
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transaction, err := h.service.CreateTransaction(input)
	if err != nil {
		response := helper.APIResponse("Failed to create transaction", 500, "error", nil)
		c.JSON(500, response)
		return
	}
	response := helper.APIResponse("Transaction created", 200, "success", transaction)
	c.JSON(200, response)
}
