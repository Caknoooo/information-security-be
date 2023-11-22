package routes

import (
	"github.com/Caknoooo/golang-clean_template/controller"
	"github.com/Caknoooo/golang-clean_template/middleware"
	"github.com/Caknoooo/golang-clean_template/services"
	"github.com/gin-gonic/gin"
)

func File(route *gin.Engine, fileController controller.FileController, jwtService services.JWTService) {
	routes := route.Group("/api/file")
	{
		routes.POST("", middleware.Authenticate(jwtService) ,fileController.UploadFile)
		routes.GET("", middleware.Authenticate(jwtService) ,fileController.GetAllFileByUser)
		routes.GET("/last", middleware.Authenticate(jwtService) ,fileController.GetLastSubmittedFilesByUserId)
		routes.GET("/get", fileController.GetFile)
	}
}