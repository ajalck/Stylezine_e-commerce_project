package handler

import (
	"ajalck/e_commerce/domain"
	services "ajalck/e_commerce/usecase/interface"
	_ "ajalck/e_commerce/utils"
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

// @title Go + Gin Workey API
// @version 1.0
// @description This is a sample server Job Portal server. You can visit the GitHub repository at https://github.com/fazilnbr/Job_Portal_Project

// @contact.name API Support
// @contact.url https://fazilnbr.github.io/mypeosolal.web.portfolio/
// @contact.email fazilkp2000@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:5050
// @BasePath
// @query.collection.format multi

// @Summary user signin
// @ID user signin
// @Tags User Authentication
// @Param userLogin body domain.User{} true "user Login"
// @Produce json
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/login [post]
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

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Token", signedToken)
		c.Writer.WriteHeader(http.StatusOK)
	}
}

// @Summary admin signin
// @ID admin signin
// @Tags admin Authentication
// @Produce json
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/registration/login [post]
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
