package interfaces

import (
	"ajalck/e_commerce/domain"
	"ajalck/e_commerce/utils"

	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	CreateUser(c *gin.Context, newUser domain.User) error
	FindUser(c *gin.Context, email string, userRole string) (domain.User, error)
	ListProducts(page, perPage int) ([]domain.ProductResponse, utils.MetaData, error)
}
