package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-finetune/app"
)

func AddRouterForFinetuneController(
	rg *gin.RouterGroup,
	fs app.FinetuneService,
) {
	ctl := FinetuneController{fs: fs}

	rg.POST("/v1/finetune", ctl.Create)
	rg.DELETE("/v1/finetune/:id", ctl.Delete)
	rg.PUT("/v1/finetune/:id", ctl.Terminate)
	rg.GET("/v1/finetune/:id/log", ctl.GetLog)
}

type FinetuneController struct {
	baseController

	fs app.FinetuneService
}

// @Summary Create
// @Description create finetune
// @Tags  Finetune
// @Param	body	body 	FinetuneCreateRequest	true	"body of creating finetune"
// @Accept json
// @Success 201 {object} app.JobInfoDTO
// @Failure 500 system_error        system error
// @Router /v1/finetune [post]
func (ctl *FinetuneController) Create(ctx *gin.Context) {
	req := FinetuneCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequest(ctx, respBadRequestBody)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequest(ctx, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, code, err := ctl.fs.Create(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, v)
	}
}

// @Summary Delete
// @Description delete finetune
// @Tags  Finetune
// @Param	id	path	string	true	"id of finetune"
// @Accept json
// @Success 204
// @Failure 500 system_error        system error
// @Router /v1/finetune/{id} [delete]
func (ctl *FinetuneController) Delete(ctx *gin.Context) {
	if err := ctl.fs.Delete(ctx.Param("id")); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfDelete(ctx)
	}
}

// @Summary Terminate
// @Description terminate finetune
// @Tags  Finetune
// @Param	id	path	string	true	"id of finetune"
// @Accept json
// @Success 202
// @Failure 500 system_error        system error
// @Router /v1/finetune/{id} [put]
func (ctl *FinetuneController) Terminate(ctx *gin.Context) {
	if err := ctl.fs.Terminate(ctx.Param("id")); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}

// @Summary GetLog
// @Description get log url of finetune for downloading
// @Tags  Finetune
// @Param	id	path	string	true	"id of finetune"
// @Accept json
// @Success 200 {object} app.LogURLDTO
// @Failure 500 system_error        system error
// @Router /v1/finetune/{id}/log [get]
func (ctl *FinetuneController) GetLog(ctx *gin.Context) {
	if v, err := ctl.fs.GetLogDownloadURL(ctx.Param("id")); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}
