package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type baseController struct{}

func (ctl baseController) sendCodeMessage(ctx *gin.Context, code string, err error) {
	data := newResponseCodeError(code, err)

	if data.Code == "" {
		logrus.Errorf("system err: %v", err)

		data.Code = errorSystemError
		ctx.JSON(http.StatusInternalServerError, data)
	} else {
		ctx.JSON(http.StatusBadRequest, data)
	}
}

func (ctl baseController) sendBadRequest(ctx *gin.Context, data responseData) {
	ctx.JSON(http.StatusBadRequest, data)
}

func (ctl baseController) sendRespOfGet(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, newResponseData(data))
}

func (ctl baseController) sendRespOfPost(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, newResponseData(data))
}

func (ctl baseController) sendRespOfPut(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusAccepted, newResponseData(data))
}

func (ctl baseController) sendRespOfDelete(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}
