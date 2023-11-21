package controller

import (
	"net/http"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/services"
	"github.com/Caknoooo/golang-clean_template/utils"
	"github.com/gin-gonic/gin"
)

type (
	PublicAccessController interface {
		PublicAccessUserFiles(ctx *gin.Context)
	}

	publicAccessController struct {
		pbs services.PublicAccessService
		js  services.JWTService
	}
)

func NewPublicAccessController(pbs services.PublicAccessService, js services.JWTService) PublicAccessController {
	return &publicAccessController{
		pbs: pbs,
		js:  js,
	}
}

func (c *publicAccessController) PublicAccessUserFiles(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	ownerId := ctx.Param("owner_id")

	var access dto.PublicAccessRequest
	if err := ctx.ShouldBind(&access); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	access.OwnerId = ownerId
	access.RequesterId = userId

	result, err := c.pbs.PublicAccessUserFiles(ctx.Request.Context(), access)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ALL_FILE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ALL_FILE, result)
	ctx.JSON(http.StatusOK, res)
}
