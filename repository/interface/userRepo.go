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
	ViewProduct(id int) (domain.Products, error)
	//Wishlist
	AddWishlist(user_id, product_id int) error
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
	//Order
	CheckOut(cart_id string, user_id, product_id, address_id int) (string, error)
	OrderSummery(order_id string) ([]domain.OrderSummery, error)
	UpdateOrder(order_id string, orderUpdates interface{}) error
}
