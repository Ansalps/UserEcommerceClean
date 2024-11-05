package controllers

import (
	"fmt"
	"net/http"

	"github.com/Ansalps/UserEcommerceClean/internal/models"
	"github.com/Ansalps/UserEcommerceClean/internal/services"
	"github.com/Ansalps/UserEcommerceClean/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.IUserService
}

func NewUserController(UserService services.IUserService) *UserController {
	return &UserController{UserService: UserService}
}

func (c *UserController) UserSignUp(ctx *gin.Context) {
	var user models.User
	err := ctx.BindJSON(&user)
	fmt.Println("error", err)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if err := utils.Validate(user); err != nil {
		fmt.Println("", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = c.UserService.UserSignUp(&user)
	if err != nil {
		if err.Error() == models.UserAlreadyExists {
			ctx.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user created",
	})
}
func (c *UserController) UserLogin(ctx *gin.Context) {
	var loginRequest models.UserLogin
	err := ctx.BindJSON(&loginRequest)
	fmt.Println("error", err)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if err := utils.Validate(loginRequest); err != nil {
		fmt.Println("", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	User, err := c.UserService.UserLogin(&loginRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	check := c.UserService.ComparePassword(loginRequest, *User)
	if !check {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "email or password is incorrect"})
		return
	}
	accessToken, _ := utils.GenerateJWT(User.Email, User.ID, "user", 1)
	//refreshToken, _ := utils.GenerateJWT(User.Email, User.ID, "user", 2)
	ctx.JSON(http.StatusOK, gin.H{"message": "Login Succesful", "user": User, "token": accessToken})
}
func (c *UserController) GetProfile(ctx *gin.Context) {
	claims, exists := ctx.Get("ID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return
	}
	// Attempt to assert claims as float64
	userIDFloat, ok := claims.(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	// Convert the float64 to a string
	userID := fmt.Sprintf("%.0f", userIDFloat)
	//var user models.User
	user, err := c.UserService.GetProfile(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user data"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	claims, exists := ctx.Get("ID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return
	}
	userID, ok := claims.(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	//userID := fmt.Sprintf("%.0f", userIDFloat)
	var updateProfileRequest models.UserUpdate
	err := ctx.BindJSON(&updateProfileRequest)
	fmt.Println("error", err)
	response := gin.H{
		"status":  false,
		"message": "failed to bind request",
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if err := utils.Validate(updateProfileRequest); err != nil {
		fmt.Println("", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     false,
			"message":    err.Error(),
			"error_code": http.StatusBadRequest,
		})
		return
	}
	err = c.UserService.UpdateProfile(uint(userID), updateProfileRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}
