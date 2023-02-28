package interfaces

import (
	"ajalck/e_commerce/domain"
	"ajalck/e_commerce/utils"

	"github.com/gin-gonic/gin"
)

type UserUseCase interface {
	CreateUser(c *gin.Context, newUser domain.User) error
	ListProducts(page, perPage int) ([]domain.ProductResponse, utils.MetaData, error)
}
type UserAuth interface {
	VerifyUser(c *gin.Context, email string, password, userRole string) (bool, error)
	FindUser(c *gin.Context, email string, role string) (domain.User, error)
}
