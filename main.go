package main

import (
	"log"
	"os"

	"github.com/Caknoooo/golang-clean_template/config"
	"github.com/Caknoooo/golang-clean_template/controller"
	"github.com/Caknoooo/golang-clean_template/middleware"
	"github.com/Caknoooo/golang-clean_template/migrations"
	"github.com/Caknoooo/golang-clean_template/repository"
	"github.com/Caknoooo/golang-clean_template/routes"
	"github.com/Caknoooo/golang-clean_template/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	var (
		db         *gorm.DB            = config.SetUpDatabaseConnection()
		jwtService services.JWTService = services.NewJWTService()

		// Repo
		userRepository         repository.UserRepository         = repository.NewUserRepository(db)
		fileRepository         repository.FileRepository         = repository.NewFileRepository(db)
		publicAccessRepository repository.PublicAccessRepository = repository.NewPublicAccessRepository(db)

		// Service
		userService         services.UserService         = services.NewUserService(userRepository, fileRepository)
		fileService         services.FileService         = services.NewFileService(fileRepository)
		publicAccessService services.PublicAccessService = services.NewPublicAccessService(publicAccessRepository, fileRepository, userRepository)

		// Controller
		userController         controller.UserController         = controller.NewUserController(userService, jwtService)
		fileController         controller.FileController         = controller.NewFileController(fileService, jwtService)
		publicAccessController controller.PublicAccessController = controller.NewPublicAccessController(publicAccessService, jwtService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())
	routes.User(server, userController, jwtService)
	routes.File(server, fileController, jwtService)
	routes.PublicAccess(server, publicAccessController, jwtService)

	if err := migrations.Seeder(db); err != nil {
		log.Fatalf("error migration seeder: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	server.Run(":" + port)
}
