package controller

import (
	"net/http"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/services"
	"github.com/Caknoooo/golang-clean_template/utils"
	"github.com/gin-gonic/gin"
)

type (
	DigitalSignatureController interface {
		CreateDigitalSignature(ctx *gin.Context)
		VerifyDigitalSignature(ctx *gin.Context)
		GetAllNotifications(ctx *gin.Context)
	}

	digitalSignatureController struct {
		digitalSignatureService services.DigitalSignatureService
	}
)

func NewDigitalSignatureController(digitalSignatureService services.DigitalSignatureService) DigitalSignatureController {
	return &digitalSignatureController{
		digitalSignatureService: digitalSignatureService,
	}
}

func (c *digitalSignatureController) CreateDigitalSignature(ctx *gin.Context) {
	var req dto.DigitalSignatureRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res, err := c.digitalSignatureService.CreateDigitalSignature(ctx, req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_DIGITAL_SIGNATURE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_DIGITAL_SIGNATURE, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *digitalSignatureController) VerifyDigitalSignature(ctx *gin.Context) {
	var req dto.VerifyDigitalSignatureRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userId := ctx.MustGet("user_id").(string)
	req.UserId = userId

	res, err := c.digitalSignatureService.VerifyDigitalSignature(ctx, req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VERIFY_DIGITAL_SIGNATURE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_VERIFY_DIGITAL_SIGNATURE, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *digitalSignatureController) GetAllNotifications(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userId := ctx.MustGet("user_id").(string)
	res, err := c.digitalSignatureService.GetAllNotifications(ctx, userId, req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_ALL_NOTIFICATIONS, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	response := utils.Response {
		Status: true,
		Message: dto.MESSAGE_SUCCESS_GET_ALL_NOTIFICATIONS,
		Data: res.Notifications,
		Meta: res.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, response)
}
