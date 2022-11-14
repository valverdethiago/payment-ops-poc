package domain

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	AccountsBasePath        = "/account"
	AccountsGetByIdPath     = "/account/:id"
	TransactionsBasePath    = "/transaction"
	TransactionsGetByIdPath = "/transaction/:id"
)

type GetAccountByIdInput struct {
	id string `uri:"id" binding:"required"`
}
type GetTransactionsByAccountIdInput struct {
	id        string `uri:"id" binding:"required"`
	accountId string `uri:"account_id" binding:"required"`
}

type AccountController interface {
	ListAllAccounts(ctx *gin.Context)
	GetAccount(ctx *gin.Context)
	ListAllTransactionsByAccount(ctx *gin.Context)
	GetTransaction(ctx *gin.Context)
	SetupRoutes(router *gin.Engine)
}

type AccountControllerImpl struct {
	accountService     AccountService
	transactionService TransactionService
}

func (controller *AccountControllerImpl) SetupRoutes(router *gin.Engine) {
	router.GET(AccountsBasePath, controller.ListAllAccounts)
	resources := router.Group(AccountsGetByIdPath)
	{
		resources.GET("/", controller.GetAccount)
		resources.GET(TransactionsBasePath, controller.ListAllTransactionsByAccount)
		resources.GET(TransactionsGetByIdPath, controller.GetTransaction)
	}
}

func NewAccountControllerImpl(accountService AccountService, transactionService TransactionService) AccountController {
	return &AccountControllerImpl{
		accountService:     accountService,
		transactionService: transactionService,
	}
}

func (controller AccountControllerImpl) ListAllAccounts(ctx *gin.Context) {
	accounts, err := controller.accountService.ListAll()
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

func (controller AccountControllerImpl) GetAccount(ctx *gin.Context) {
	payload, err := parseGetByIdInput(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ID, err := parseUUID(payload.id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := controller.accountService.GetAccountSnapshot(ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, account)
}
func (controller AccountControllerImpl) ListAllTransactionsByAccount(ctx *gin.Context) {
	payload, err := parseGetByIdInput(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountId, err := parseUUID(payload.id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transactions, err := controller.transactionService.FindAllTransactionsByAccount(accountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if transactions == nil || len(*transactions) < 1 {
		ctx.Status(http.StatusNoContent)
		return
	}
	ctx.JSON(http.StatusOK, transactions)
}

func (controller AccountControllerImpl) GetTransaction(ctx *gin.Context) {
	payload, err := parseGetTransactionsByAccountIdInput(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountId, err := parseUUID(payload.accountId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transactionId, err := parseUUID(payload.id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction, err := controller.transactionService.FindByAccountIdAndTransactionId(accountId, transactionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if transaction == nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, transaction)
}

func parseGetTransactionsByAccountIdInput(ctx *gin.Context) (*GetTransactionsByAccountIdInput, error) {
	return &GetTransactionsByAccountIdInput{
		id:        ctx.Param("id"),
		accountId: ctx.Param("account_id"),
	}, nil
}

func parseGetByIdInput(ctx *gin.Context) (*GetAccountByIdInput, error) {
	return &GetAccountByIdInput{
		ctx.Param("id"),
	}, nil
}
