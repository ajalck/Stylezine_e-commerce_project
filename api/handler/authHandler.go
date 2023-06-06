package handler

import (
	"ajalck/e_commerce/domain"
	services "ajalck/e_commerce/usecase/interface"
	"ajalck/e_commerce/utils"
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

type Signin struct {
	Username string `json:"username" gorm:"not null" binding:"required,email"`
	Password string `json:"password" gorm:"not null" binding:"required,min=6"`
}

// @Summary user signin
// @ID user signin
// @Tags 11.User Authentication
// @Param userLogin body Signin{} true "user Login"
// @Produce json
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/login [post]
func (uh *UserAuthHandler) UserSignin(c *gin.Context) {
	data := &Signin{}
	if err := c.Bind(&data); err != nil {
		response := utils.ErrorResponse("Authentication failed !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusExpectationFailed)
		utils.ResponseJSON(c, response)
		return
	}
	userLogin := &domain.User{
		Email:     data.Username,
		Password:  data.Password,
		User_Role: "user",
	}
	isVerified, _ := uh.userAuth.VerifyUser(userLogin.Email, userLogin.Password, userLogin.User_Role)
	if !isVerified {
		response := utils.ErrorResponse("Authentication failed !", "", nil)
		c.Writer.WriteHeader(http.StatusExpectationFailed)
		utils.ResponseJSON(c, response)
	} else {
		user, _ := uh.userAuth.FindUser(userLogin.Email, userLogin.User_Role)

		signedToken := uh.jwtService.GenerateToken(user.User_ID, user.Email, "user")

		user.Password = ""

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Token", signedToken)
		response := utils.SuccessResponse("Logged in successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// @Summary admin signin
// @ID admin signin
// @Tags 2.admin Authentication
// @Produce json
// @Param adminLogin body Signin{} true "admin login"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/registration/login [post]
func (ah *AdminAuthHandler) AdminSignin(c *gin.Context) {
	data := &Signin{}
	if err := c.Bind(&data); err != nil {
		response := utils.ErrorResponse("Authentication failed !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusExpectationFailed)
		utils.ResponseJSON(c, response)
		return
	}
	adminLogin := &domain.User{
		Email:     data.Username,
		Password:  data.Password,
		User_Role: "admin",
	}
	isVerified, _ := ah.adminAuth.VerifyAdmin(adminLogin.Email, adminLogin.Password, adminLogin.User_Role)
	if !isVerified {
		response := utils.ErrorResponse("Authentication failed!", "", nil)
		c.Writer.WriteHeader(http.StatusExpectationFailed)
		utils.ResponseJSON(c, response)
	} else {
		admin, _ := ah.adminAuth.FindAdmin(adminLogin.Email, adminLogin.User_Role)

		signedToken := ah.jwtService.GenerateToken(admin.User_ID, admin.Email, "admin")

		admin.Password = ""

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Token", signedToken)
		response := utils.SuccessResponse("Logged in successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}
