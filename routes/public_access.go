package routes

import (
	"github.com/Caknoooo/golang-clean_template/controller"
	"github.com/Caknoooo/golang-clean_template/middleware"
	"github.com/Caknoooo/golang-clean_template/services"
	"github.com/gin-gonic/gin"
)

func PublicAccess(route *gin.Engine, pubAccessController controller.PublicAccessController, jwtService services.JWTService) {
	routes := route.Group("/api/public_access")
	{
		routes.GET("/:owner_id", middleware.Authenticate(jwtService), pubAccessController.PublicAccessUserFiles)
	}
}
