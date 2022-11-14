package domain

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
)

const (
	syncRequestPathById = "/request/:id"
	newSyncRequestPath  = "/request/:accountId/:syncType"
)

type requestId struct {
	ID string `uri:"id" binding: "required"`
}

type newSyncRequest struct {
	AccountId uuid.UUID `uri:"accountId" binding: "required"`
	SyncType  SyncType  `uri:"syncType" binding: "required,enum"`
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
	var req newSyncRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	created, err := controller.service.Request(req.AccountId, req.SyncType)
	if err != nil {
		if err == ErrorInvalidValueForSyncType {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
	}
	ctx.JSON(http.StatusOK, created)
}

func (controller *SyncRequestController) assertRequestExists(ctx *gin.Context) (*SyncRequest, error) {
	req, err := parseIdRequest(ctx)
	if err != nil {
		return nil, err
	}
	syncRequest, err := controller.service.Find(req.ID)
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
