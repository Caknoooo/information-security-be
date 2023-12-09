package main

import (
	"os"

	"github.com/Caknoooo/golang-clean_template/config"
	"github.com/Caknoooo/golang-clean_template/controller"
	"github.com/Caknoooo/golang-clean_template/middleware"

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
		userRepository             repository.UserRepository             = repository.NewUserRepository(db)
		fileRepository             repository.FileRepository             = repository.NewFileRepository(db)
		privateAccessRepository    repository.PrivateAccessRepository    = repository.NewPrivateAccessRepository(db)
		digitalSignatureRepository repository.DigitalSignatureRepository = repository.NewDigitalSignatureRepository(db)

		// Service
		userService             services.UserService             = services.NewUserService(userRepository, fileRepository)
		fileService             services.FileService             = services.NewFileService(fileRepository)
		privateAccessService    services.PrivateAccessService    = services.NewPrivateAccessService(userRepository, privateAccessRepository, fileRepository)
		digitalSignatureService services.DigitalSignatureService = services.NewDigitalSignatureService(digitalSignatureRepository, userRepository)

		// Controller
		userController             controller.UserController             = controller.NewUserController(userService, jwtService)
		fileController             controller.FileController             = controller.NewFileController(fileService, jwtService)
		privateAccessController    controller.PrivateAccessController    = controller.NewPrivateAccessController(privateAccessService)
		digitalSignatureController controller.DigitalSignatureController = controller.NewDigitalSignatureController(digitalSignatureService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())
	routes.User(server, userController, jwtService)
	routes.File(server, fileController, jwtService)
	routes.PrivateAccess(server, privateAccessController, jwtService)
	routes.DigitalSignature(server, digitalSignatureController, jwtService)

	// if err := migrations.Seeder(db); err != nil {
	// 	log.Fatalf("error migration seeder: %v", err)
	// }

	server.Static("/storage", "./storage")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	server.Run(":" + port)
}