package interfaces

import (
	"ajalck/e_commerce/domain"
	"ajalck/e_commerce/utils"
)

type UserUseCase interface {
	CreateUser(newUser domain.User) error
	ListProducts(page, perPage int) ([]domain.ProductResponse, utils.MetaData, error)
	//Wishlist
	AddWishlist(user_id, Product_id string) error
	ViewWishList(user_id string, page, perPage int) ([]domain.WishListResponse, utils.MetaData, error)
	DeleteWishList(wishlist_id string) error
	//Cart
	AddCart(user_id, product_id string) (error, string)
	ViewCart(user_id string, page, perPage int) ([]domain.CartResponse, float32, utils.MetaData, error)
	DeleteCart(user_id, product_id string) error
	//Coupon
	ListCoupon(user_id, product_id string) ([]domain.CouponResponse, error)
	ApplyCoupon(cart_id, order_id, coupon_id string) error
	CancelCoupon(cart_id, order_id, coupon_id string) error
	//Shipping
	AddShippingDetails(user_id string, newAddress domain.ShippingDetails) error
	ListShippingDetails(user_id string) ([]domain.ShippingDetailsResponse, error)
	DeleteShippingDetails(user_id, address_id string) error
	//order
	CheckOut(cart_id, user_id, product_id, address_id string) (string, error)
	OrderSummery(user_id string) (interface{}, domain.OrderSummery, error)
	CancelOrder(order_id string) error
}
type UserAuth interface {
	VerifyUser(email string, password, userRole string) (bool, error)
	FindUser(email string, role string) (domain.User, error)
}
