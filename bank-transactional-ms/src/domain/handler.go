package domain

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	testPathById = "/test"
)

type TestController struct {
}

func NewTestController() *TestController {
	return &TestController{}
}

func (controller *TestController) SetupRoutes(router *gin.Engine) {
	router.GET(testPathById, controller.Test)
}

func (controller *TestController) Test(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "{ 'message' : 'Hello stranger' }")
}
