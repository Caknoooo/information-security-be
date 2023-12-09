package routes

import (
	"github.com/Caknoooo/golang-clean_template/controller"
	"github.com/Caknoooo/golang-clean_template/middleware"
	"github.com/Caknoooo/golang-clean_template/services"

	"github.com/gin-gonic/gin"
)

func DigitalSignature(routes *gin.Engine, digitalSignatureController controller.DigitalSignatureController, jwtService services.JWTService) {
	digitalSignature := routes.Group("/api/digital_signature")
	{
		digitalSignature.POST("", middleware.Authenticate(jwtService), digitalSignatureController.CreateDigitalSignature)
		digitalSignature.POST("/verify", middleware.Authenticate(jwtService), digitalSignatureController.VerifyDigitalSignature)
		digitalSignature.GET("/notifications", middleware.Authenticate(jwtService), digitalSignatureController.GetAllNotifications)
	}
}
