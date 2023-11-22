package controller

import (
	"net/http"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/services"
	"github.com/Caknoooo/golang-clean_template/utils"
	"github.com/gin-gonic/gin"
)

type (
	PrivateAccessController interface {
		Create(ctx *gin.Context)
		GetAllPrivateAccessRequestByUserId(ctx *gin.Context)
		GetAllPrivateAccessOwnerByUserId(ctx *gin.Context)
		UpdatePrivateAccess(ctx *gin.Context)
		SendEncryptionKey(ctx *gin.Context)
	}

	privateAccessController struct {
		privateAccessService services.PrivateAccessService
	}
)

func NewPrivateAccessController(privateAccessService services.PrivateAccessService) PrivateAccessController {
	return &privateAccessController{
		privateAccessService: privateAccessService,
	}
}

func (c *privateAccessController) Create(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	var req dto.PrivateAccessRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	req.UserReqId = userId

	result, err := c.privateAccessService.Create(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_PRIVATE_ACCESS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_PRIVATE_ACCESS, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *privateAccessController) GetAllPrivateAccessRequestByUserId(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	result, err := c.privateAccessService.GetAllPrivateAccessRequestByUserId(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ALL_PRIVATE_ACCESS_REQUEST, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ALL_PRIVATE_ACCESS_REQUEST, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *privateAccessController) GetAllPrivateAccessOwnerByUserId(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)

	result, err := c.privateAccessService.GetAllPrivateAccessOwnerByUserId(ctx.Request.Context(), userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ALL_PRIVATE_ACCESS_OWNER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_ALL_PRIVATE_ACCESS_OWNER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *privateAccessController) UpdatePrivateAccess(ctx *gin.Context) {
	var req dto.UpdatePrivateAccessRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userId := ctx.MustGet("user_id").(string)
	id := ctx.Param("id")
	req.ID = id

	result, err := c.privateAccessService.UpdatePrivateAccess(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_PRIVATE_ACCESS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_PRIVATE_ACCESS, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *privateAccessController) SendEncryptionKey(ctx *gin.Context) {
	var req dto.SendEncryptionKeyRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userReqId := ctx.MustGet("user_id").(string)
	ownerId := ctx.Param("owner_id")
	req.OwnerId = ownerId

	result, err := c.privateAccessService.SendEncryptionKey(ctx.Request.Context(), req, userReqId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_SEND_ENCRYPTION_KEY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_SEND_ENCRYPTION_KEY, result)
	ctx.JSON(http.StatusOK, res)
}