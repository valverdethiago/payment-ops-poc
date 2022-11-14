package domain

import (
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/restclient"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	webhookPath = "webhook/:type"
)

type WebHookInput struct {
	SyncType SyncType `uri:"type" binding:"required"`
}

type WebHookController struct {
	syncRequestService SyncRequestService
	balanceService     BalanceService
	transactionService TransactionService
	accountRepository  AccountRepository
}

func NewWebHookControllerImpl(syncRequestService SyncRequestService, balanceService BalanceService,
	transactionService TransactionService, accountRepository AccountRepository) *WebHookController {
	return &WebHookController{
		syncRequestService: syncRequestService,
		balanceService:     balanceService,
		transactionService: transactionService,
		accountRepository:  accountRepository,
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
		controller.handleBalancesCallback(ctx)
	case SyncTypeTransactions:
		controller.handleTransactionsCallback(ctx)
	}
}

func (controller *WebHookController) handleBalancesCallback(ctx *gin.Context) {
	payload, err := ParseBalancesWebHookPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := controller.findAccount(payload.Event.AccountID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if payload.Error != nil && payload.Error.Message != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": payload.Error.Message})
		controller.syncRequestService.ChangeToFailingStatus(account.InternalAccountId, SyncTypeBalances,
			payload.Error.Message)
	}
	err = controller.balanceService.UpdateAccountBalance(payload.Event.AccountID, payload.Data.Amount.Amount,
		payload.Data.Amount.Currency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	controller.syncRequestService.ChangeToSuccessfulStatus(account.InternalAccountId, SyncTypeBalances)
	ctx.Status(http.StatusOK)

}

func (controller *WebHookController) handleTransactionsCallback(ctx *gin.Context) {
	payload, err := ParseTransactionsWebHookPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := controller.findAccount(payload.Event.AccountID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if payload.Error != nil && payload.Error.Message != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": payload.Error.Message})
		controller.syncRequestService.ChangeToFailingStatus(account.InternalAccountId, SyncTypeTransactions,
			payload.Error.Message)
	}
	err = controller.transactionService.UpdateAccountTransactions(payload.Event.AccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

func (controller *WebHookController) findAccount(id string) (*Account, error) {
	return controller.accountRepository.FindByProviderAccountId(id)
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
