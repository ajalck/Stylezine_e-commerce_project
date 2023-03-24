package interfaces

import (
	"ajalck/e_commerce/domain"
	"ajalck/e_commerce/utils"
)

type AdminRepository interface {
	CreateAdmin(admin domain.User) error
	FindAdmin(email string, userRole string) (domain.User, error)

	ListUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error)
	ListBlockedUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error)
	ListActiveUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error)
	ViewUser(id string) (domain.UserResponse, error)
	BlockUser(id string)
	UnblockUser(id string)

	AddProducts(products domain.Products) (string, error)
	ListProducts(page, perPage int) ([]domain.ProductResponse, utils.MetaData, error)
	EditProducts(products domain.Products) error
	DeleteProducts(product_id string) error

	AddCategory(category domain.Category) error
	ListCategory() ([]domain.Category, error)
	ViewCategory(category domain.Category) (domain.Category, error)
	EditCategory(category domain.Category) error
	DeleteCategory(category domain.Category) error

	AddBrand(brand domain.Brand) error
	ListBrands() ([]domain.Brand, error)
	ViewBrand(brand_id uint) (domain.Brand, error)
	EditBrand(brand domain.Brand) error
	DeleteBrand(brand domain.Brand) error

	AddCoupon(coupon domain.Coupon) error
	ListCoupon(page, perPage int) ([]domain.CouponResponse, utils.MetaData, error)
	DeleteCoupon(coupon_id string) error

	SalesReport(page, perPage int) (interface{}, utils.MetaData, error)
}
