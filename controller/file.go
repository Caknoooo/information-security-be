package controller

import (
	"net/http"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/services"
	"github.com/Caknoooo/golang-clean_template/utils"
	"github.com/gin-gonic/gin"
)

type (
	FileController interface {
		UploadFile(ctx *gin.Context)
	}

	fileController struct {
		fileService services.FileService
		jwtService services.JWTService
	}
)

func NewFileController(fileService services.FileService, jwtService services.JWTService) FileController {
	return &fileController{
		fileService: fileService,
		jwtService: jwtService,
	}
}

func (c *fileController) UploadFile(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userId, err := c.jwtService.GetUserIDByToken(token)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER_TOKEN, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.UploadFileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.fileService.UploadFile(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPLOAD_FILE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPLOAD_FILE, result)
	ctx.JSON(http.StatusOK, res)
}