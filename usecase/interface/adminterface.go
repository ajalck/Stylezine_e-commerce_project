package interfaces

import (
	"ajalck/e_commerce/domain"
	"ajalck/e_commerce/utils"

	"github.com/gin-gonic/gin"
)

type AdminUseCase interface {
	CreateAdmin(c *gin.Context, admin domain.User) error

	ListUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error)
	ListBlockedUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error)
	ListActiveUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error)
	ViewUser(id int) (domain.UserResponse, error)
	BlockUser(id int) error
	UnblockUser(id int) error

	AddCategory(category domain.Category) error
	EditCategory(category domain.Category) error
	DeleteCategory(category domain.Category) error

	AddBrand(brand domain.Brand) error
	EditBrand(brand domain.Brand) error
	DeleteBrand(brand domain.Brand) error

	AddProducts(products domain.Products) error
	EditProducts(products domain.Products) error
	DeleteProducts(products domain.Products) error
}
type AdminAuth interface {
	VerifyAdmin(c *gin.Context, email string, password, userRole string) (bool, error)
}
