package routes

import (
	"github.com/Caknoooo/golang-clean_template/controller"
	"github.com/Caknoooo/golang-clean_template/middleware"
	"github.com/Caknoooo/golang-clean_template/services"

	"github.com/gin-gonic/gin"
)

func PrivateAccess(route *gin.Engine, privateAccessController controller.PrivateAccessController, jwtService services.JWTService) {
	routes := route.Group("/api/private-access", middleware.Authenticate(jwtService))
	{
		routes.POST("", middleware.Authenticate(jwtService), privateAccessController.Create)
		routes.GET("/request", middleware.Authenticate(jwtService), privateAccessController.GetAllPrivateAccessRequestByUserId)
		routes.GET("/owner", middleware.Authenticate(jwtService), privateAccessController.GetAllPrivateAccessOwnerByUserId)
		routes.PATCH("/:id", middleware.Authenticate(jwtService), privateAccessController.UpdatePrivateAccess)
		routes.DELETE("/:id", middleware.Authenticate(jwtService), privateAccessController.DeletePrivateAccess)
		routes.POST("/send-key/:owner_id", middleware.Authenticate(jwtService), privateAccessController.SendEncryptionKey)
	}
}
