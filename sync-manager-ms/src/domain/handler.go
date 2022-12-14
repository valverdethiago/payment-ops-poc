package domain

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

const (
	syncRequestPathById = "/request/:id"
	newSyncRequestPath  = "/request/:accountId/:syncType"
)

type requestId struct {
	ID string `uri:"id" binding:"required"`
}

type newSyncRequest struct {
	AccountId string `uri:"accountId" binding:"required"`
	SyncType  string `uri:"syncType" binding:"required,enum"`
}

type SyncRequestController struct {
	service SyncRequestService
}

func NewSyncRequestController(service SyncRequestService) *SyncRequestController {
	return &SyncRequestController{
		service: service,
	}
}

func (controller *SyncRequestController) SetupRoutes(router *gin.Engine) {
	router.GET(syncRequestPathById, controller.FindById)
	router.POST(newSyncRequestPath, controller.NewRequest)
}

func (controller *SyncRequestController) FindById(ctx *gin.Context) {
	req, err := controller.assertRequestExists(ctx)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, req)
}

func (controller *SyncRequestController) NewRequest(ctx *gin.Context) {
	req, err := parseNewRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	AccountId, err := ParseUUID(req.AccountId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	SyncType, err := ParseSyncType(req.SyncType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	now := time.Now()
	created, err := controller.service.Request(AccountId, SyncType, now)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if created.CreatedAt.Before(now) {
		ctx.JSON(http.StatusAlreadyReported, created)
		return
	}
	ctx.JSON(http.StatusCreated, created)
}

func parseNewRequest(ctx *gin.Context) (*newSyncRequest, error) {
	return &newSyncRequest{
		AccountId: ctx.Param("accountId"),
		SyncType:  ctx.Param("syncType"),
	}, nil
}

func (controller *SyncRequestController) assertRequestExists(ctx *gin.Context) (*SyncRequest, error) {
	req, err := parseIdRequest(ctx)
	if err != nil {
		return nil, err
	}
	ID := ParseBson(req.ID)
	syncRequest, err := controller.service.Find(ID)
	if err != nil {
		if err == mgo.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "No request found for the given id"})
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
	}
	return syncRequest, err
}

func parseIdRequest(ctx *gin.Context) (requestId, error) {
	var req requestId
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	return req, err
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
