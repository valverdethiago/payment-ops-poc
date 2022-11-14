package domain

import (
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/restclient"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const (
	webhookPath = "webhook/:type"
)

type WebHookInput struct {
	SyncType SyncType `uri:"type" binding:"required"`
}

type WebHookController struct {
	SyncRequestService SyncRequestService
	BalanceService     BalanceService
	TransactionService TransactionService
}

func NewWebHookControllerImpl(SyncRequestService SyncRequestService, BalanceService BalanceService,
	TransactionService TransactionService) *WebHookController {
	return &WebHookController{
		SyncRequestService: SyncRequestService,
		BalanceService:     BalanceService,
		TransactionService: TransactionService,
	}
}

func (controller *WebHookController) SetupRoutes(router *gin.Engine) {
	router.POST(webhookPath, controller.WebhookCall)
}

func (controller *WebHookController) WebhookCall(ctx *gin.Context) {
	input, err := ParseWebHookInput(ctx)
	if err != nil {
		return
	}
	switch input.SyncType {
	case SyncTypeBalances:
		controller.HandleBalancesCallback(ctx)
	case SyncTypeTransactions:
		controller.HandleTransactionsCallback(ctx)
	}
	ctx.Status(http.StatusOK)
}

func (controller *WebHookController) HandleBalancesCallback(ctx *gin.Context) {
	payload, err := ParseBalancesWebHookPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Error != nil && payload.Error.Message != nil {
		controller.SyncRequestService.ChangeToFailingStatus(bson.ObjectIdHex(payload.Event.AccountID),
			payload.Error.Message)
	}
	controller.BalanceService.UpdateAccountBalance(payload.Event.AccountID, payload.Data.Amount.Amount,
		payload.Data.Amount.Currency)

}

func (controller *WebHookController) HandleTransactionsCallback(ctx *gin.Context) {
	payload, err := ParseTransactionsWebHookPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Error != nil && payload.Error.Message != nil {
		controller.SyncRequestService.ChangeToFailingStatus(bson.ObjectIdHex(payload.Event.AccountID), payload.Error.Message)
	}
	controller.TransactionService.UpdateAccountTransactions(payload.Event.AccountID)
}

func ParseWebHookInput(ctx *gin.Context) (*WebHookInput, error) {
	var req WebHookInput
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}
	return &req, nil
}

func ParseBalancesWebHookPayload(ctx *gin.Context) (*restclient.BalancesWebHookPayload, error) {
	var req restclient.BalancesWebHookPayload
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}
	return &req, err
}

func ParseTransactionsWebHookPayload(ctx *gin.Context) (*restclient.TransactionsWebHookPayload, error) {
	var req restclient.TransactionsWebHookPayload
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}
	return &req, err
}
