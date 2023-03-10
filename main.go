package main

import (
	"ajalck/e_commerce/api/handler"
	"ajalck/e_commerce/api/middleware"
	"ajalck/e_commerce/api/servers"
	"ajalck/e_commerce/config"
	_ "ajalck/e_commerce/docs"
	"ajalck/e_commerce/repository"
	repoInt "ajalck/e_commerce/repository/interface"
	"ajalck/e_commerce/usecase"
	services "ajalck/e_commerce/usecase/interface"
	_ "ajalck/e_commerce/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Failed to connect environment variables")
	}
}

func main() {
	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())
	// Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	db := config.ConnectDB()
	config.SyncDB(db)

	var (
		userRepo        repoInt.UserRepository    = repository.NewUserRepository(db)
		userUseCase     services.UserUseCase      = usecase.NewUserUseCase(userRepo)
		userHandler     *handler.UserHandler      = handler.NewUserHandler(userUseCase)
		userAuth        services.UserAuth         = usecase.NewUserAuthService(userRepo)
		jwtService      services.JwtServices      = usecase.NewJWTService()
		userAuthHandler *handler.UserAuthHandler  = handler.NewUserAuthHandler(userAuth, jwtService)
		userMiddleware  middleware.UserMiddleware = middleware.NewUserMiddleware(jwtService)
	)
	var (
		adminRepo        repoInt.AdminRepository   = repository.NewAdminRepository(db)
		adminUseCase     services.AdminUseCase     = usecase.NewAdminUseCase(adminRepo)
		adminHandler     *handler.AdminHandler     = handler.NewAdminHandler(adminUseCase)
		adminAuth        services.AdminAuth        = usecase.NewAdminAuthService(adminRepo)
		adminAuthHandler *handler.AdminAuthHandler = handler.NewAdminAuthHandler(adminAuth, jwtService)
	)
	//routing
	servers.UserServer(router, *userHandler, *userAuthHandler, userMiddleware)
	servers.AdminServer(router, *adminHandler, *adminAuthHandler)
	router.Run(":" + port)
}
