package handler

import (
	"ajalck/e_commerce/domain"
	services "ajalck/e_commerce/usecase/interface"
	"ajalck/e_commerce/utils"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

func NewUserHandler(userUseCase services.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	var newUser domain.User
	if err := c.Bind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Signup Inputs"})
		return
	}

	err := uh.userUseCase.CreateUser(c, newUser)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't Create a new user"})
		return
	} else {
		c.JSON(http.StatusFound, gin.H{"message": "New user created successfully"})
	}
}
func (uh *UserHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	users, metaData, err := uh.userUseCase.ListProducts(page, perPage)
	type Page struct {
		users    []domain.ProductResponse
		metaData utils.MetaData
	}
	result := Page{
		users:    users,
		metaData: metaData,
	}
	if err == nil {
		c.JSON(200, result.users)
		c.JSON(http.StatusFound, result.metaData)
		return
	}
	if err != nil {
		c.JSON(400, err.Error())
	}
}
