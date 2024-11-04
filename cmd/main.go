package main

import (
	"github.com/Ansalps/UserEcommerceClean/internal/controllers"
	"github.com/Ansalps/UserEcommerceClean/internal/database"
	"github.com/Ansalps/UserEcommerceClean/internal/middleware"
	"github.com/Ansalps/UserEcommerceClean/internal/repository"
	"github.com/Ansalps/UserEcommerceClean/internal/services"
	"github.com/gin-gonic/gin"
)

func init() {
	database.Initialize()
	database.AutoMigrate()
}
func main() {
	router := gin.Default()
	userRepo := repository.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)
	//User Routes
	//router.POST("storename", userController.StoreName)

	router.POST("user-signup", userController.UserSignUp)
	router.POST("user-login", userController.UserLogin)
	userGroup := router.Group("user/")
	userGroup.Use(middleware.JWTMIddleware("user"))
	userGroup.GET("profile", userController.GetProfile)
	userGroup.PUT("profile", userController.UpdateProfile)
	//router.RegisterUrls(router)
	//router.LoadHTMLGlob("templates/*")
	router.Run(":5000")
}
