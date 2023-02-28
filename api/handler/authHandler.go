package handler

import (
	"ajalck/e_commerce/domain"
	services "ajalck/e_commerce/usecase/interface"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAuthHandler struct {
	userAuth   services.UserAuth
	jwtService services.JwtServices
}
type AdminAuthHandler struct {
	adminAuth  services.AdminAuth
	jwtService services.JwtServices
}

func NewUserAuthHandler(
	userAuth services.UserAuth,
	jwtService services.JwtServices,
) *UserAuthHandler {
	return &UserAuthHandler{
		userAuth:   userAuth,
		jwtService: jwtService}
}

func NewAdminAuthHandler(
	adminAuth services.AdminAuth,
	jwtService services.JwtServices,
) *AdminAuthHandler {
	return &AdminAuthHandler{
		adminAuth:  adminAuth,
		jwtService: jwtService,
	}
}

func (uh *UserAuthHandler) UserSignin(c *gin.Context) {
	var userLogin domain.User
	if err := c.Bind(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login inputs"})
		return
	}
	userLogin.User_Role = "user"
	isVerified, _ := uh.userAuth.VerifyUser(c, userLogin.Email, userLogin.Password, userLogin.User_Role)
	if !isVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs"})
	} else {
		user, _ := uh.userAuth.FindUser(c, userLogin.Email, userLogin.User_Role)

		signedToken := uh.jwtService.GenerateToken(int(user.ID), user.Email, "user")

		user.Password = ""
		user.Token = signedToken

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Token", user.Token)
		c.Writer.WriteHeader(http.StatusOK)
	}
}
func (ah *AdminAuthHandler) AdminSignin(c *gin.Context) {
	var adminLogin domain.User
	if err := c.Bind(&adminLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login inputs"})
		return
	}
	adminLogin.User_Role = "admin"
	isVerified, _ := ah.adminAuth.VerifyAdmin(c, adminLogin.Email, adminLogin.Password, adminLogin.User_Role)
	if !isVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs"})
	} else {
		admin := adminLogin

		signedToken := ah.jwtService.GenerateToken(int(admin.ID), admin.Email, "user")

		fmt.Print("\n\n", "signed token:", signedToken, "\n\n")
	}
}
