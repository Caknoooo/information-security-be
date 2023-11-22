package controller

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/services"
	"github.com/Caknoooo/golang-clean_template/utils"
	"github.com/gin-gonic/gin"
)

type (
	FileController interface {
		UploadFile(ctx *gin.Context)
		GetAllFileByUser(ctx *gin.Context)
		GetLastSubmittedFilesByUserId(ctx *gin.Context)
		GetFile(ctx *gin.Context)
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
	userId := ctx.MustGet("user_id").(string)
	mode := ctx.Query("mode")

	var req dto.UploadFileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.fileService.UploadFile(ctx.Request.Context(), req, userId, mode)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPLOAD_FILE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPLOAD_FILE, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *fileController) GetLastSubmittedFilesByUserId(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	result, err := c.fileService.GetLastSubmittedFilesByUserId(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LAST_SUBMITTED_FILES, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LAST_SUBMITTED_FILES, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *fileController) GetAllFileByUser(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	result, err := c.fileService.GetAllFileByUser(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ALL_FILE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ALL_FILE, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *fileController) GetFile(ctx *gin.Context) {
	mode := ctx.Query("mode")
	filename := ctx.Query("filename")

	fileDecrypt, err := c.fileService.DecryptFileData(ctx.Request.Context(), filename, mode)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_FILE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	fmt.Println(mode, filename, fileDecrypt)

	data := strings.Split(fileDecrypt, "/")
	filePath := utils.PATH + "/" + data[0] + "/" + data[1] + "/" + data[2]

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_FILE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	ctx.File(filePath)
}
