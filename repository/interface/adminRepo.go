package interfaces

import (
	"ajalck/e_commerce/domain"
	"ajalck/e_commerce/utils"

	"github.com/gin-gonic/gin"
)

type AdminRepository interface {
	CreateAdmin(c *gin.Context, admin domain.User) error
	FindAdmin(c *gin.Context, email string, userRole string) (domain.User, error)

	ListUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error)
	ListBlockedUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error)
	ListActiveUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error)
	ViewUser(id int) (domain.UserResponse, error)
	BlockUser(id int)
	UnblockUser(id int)

	AddProducts(products domain.Products) (string,error)
	EditProducts(products domain.Products) error
	DeleteProducts(products domain.Products) error

	ViewCategory(category domain.Category) (domain.Category, error)
	EditCategory(category domain.Category) error
	AddCategory(category domain.Category) error
	DeleteCategory(category domain.Category) error

	ViewBrand(brand_id uint) (domain.Brand, error)
	AddBrand(brand domain.Brand) error
	EditBrand(brand domain.Brand) error
	DeleteBrand(brand domain.Brand) error

	AddCoupon(coupon domain.Coupon)error
	ListCoupon(page, perPage int) ([]domain.CouponResponse, utils.MetaData, error)
	DeleteCoupon(coupon_id int) error 
}
