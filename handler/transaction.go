package handler

import (
	"net/http"
	"rocketship/helper"
	"rocketship/transaction"
	"rocketship/user"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (handler *transactionHandler) FindTransactionByCampaignID(context *gin.Context) {
	var input transaction.FindTransactionByIDInput

	err := context.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse(
			"Failed to get transactions with this Campaign ID",
			http.StatusUnprocessableEntity,
			"error",
			err,
		)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := context.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactionList, err := handler.service.FindTransactionByCampaignID(input)
	if err != nil {
		response := helper.APIResponse(
			"Failed to get transactions due to server error",
			http.StatusBadRequest,
			"error",
			err.Error(),
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse(
		"Transaction fetched",
		http.StatusOK,
		"success",
		transaction.FormatCampaignTransactionList(transactionList),
	)
	context.JSON(http.StatusOK, response)
}

func (handler *transactionHandler) FindTransactionByUserID(context *gin.Context) {
	currentUser := context.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactionList, err := handler.service.FindTransactionByUserID(userID)
	if err != nil {
		response := helper.APIResponse(
			"Failed to get transactions due to server error",
			http.StatusBadRequest,
			"error",
			err.Error(),
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse(
		"Transaction fetched",
		http.StatusOK,
		"success",
		transaction.FormatUserTransactionList(transactionList),
	)
	context.JSON(http.StatusOK, response)
}

func (handler *transactionHandler) CreateTransaction(context *gin.Context) {
	var input transaction.CreateTransactionInput

	err := context.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse(
			"Failed to create transaction due to bad inputs",
			http.StatusUnprocessableEntity,
			"failed",
			err.Error(),
		)
		context.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := context.MustGet("currentUser").(user.User)
	input.User = currentUser

	newTransaction, err := handler.service.CreateTransaction(input)
	if err != nil {
		response := helper.APIResponse(
			"Failed to create transaction due to server error",
			http.StatusBadRequest,
			"failed",
			err.Error(),
		)
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse(
		"Transaction successfully created!",
		http.StatusOK,
		"success",
		newTransaction,
	)
	context.JSON(http.StatusOK, response)
}