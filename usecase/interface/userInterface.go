package interfaces

import (
	"ajalck/e_commerce/domain"
	"ajalck/e_commerce/utils"

	"github.com/gin-gonic/gin"
)

type UserUseCase interface {
	CreateUser(c *gin.Context, newUser domain.User) error
	ListProducts(page, perPage int) ([]domain.ProductResponse, utils.MetaData, error)
	//Wishlist
	AddWishlist(user_id, Product_id int) error
	ViewWishList(user_id, page, perPage int) ([]domain.WishListResponse, utils.MetaData, error)
	DeleteWishList(user_id, product_id int) error
	//Cart
	AddCart(user_id, product_id int) (error, string)
	ViewCart(user_id, page, perPage int) ([]domain.CartResponse, utils.MetaData, error)
	DeleteCart(user_id, product_id int) error
	//Coupon
	ListCoupon(user_id, product_id int) ([]domain.CouponResponse, error)
	ApplyCoupon(cart_id, order_id string, coupon_id int) error
	CancelCoupon(cart_id, order_id string, coupon_id int) error
	//Shipping
	AddShippingDetails(user_id int, newAddress domain.ShippingDetails) error
	ListShippingDetails(user_id int) ([]domain.ShippingDetailsResponse, error)
	DeleteShippingDetails(user_id, address_id int) error
	//order
	CheckOut(cart_id string, user_id, product_id, address_id int) (string, error)
	OrderSummery(order_id string) ([]domain.OrderSummery, error)
	UpdateOrder(orders_id string) error
}
type UserAuth interface {
	VerifyUser(c *gin.Context, email string, password, userRole string) (bool, error)
	FindUser(c *gin.Context, email string, role string) (domain.User, error)
}
