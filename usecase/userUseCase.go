package usecase

import (
	"ajalck/e_commerce/domain"
	repoInt "ajalck/e_commerce/repository/interface"
	services "ajalck/e_commerce/usecase/interface"
	"ajalck/e_commerce/utils"
	"crypto/md5"
	"errors"
	"fmt"
)

type userUseCase struct {
	userRepo repoInt.UserRepository
}

func NewUserUseCase(repo repoInt.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}
func (uc *userUseCase) CreateUser(newUser domain.User) error {

	newUser.User_Role = "user"

	if _, err := uc.userRepo.FindUser(newUser.Email, newUser.User_Role); err == nil {

		return errors.New("user already exists")
	}

	//hashing password

	newUser.Password = HashPassword(newUser.Password)
	newUser.User_ID = utils.GenerateID()
	newUser.Status = "active"
	newUser.Level = "bronze"
	err := uc.userRepo.CreateUser(newUser)
	if err != nil {
		return err
	}
	return nil

}

func HashPassword(password string) string {
	bytes := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	return string(bytes)
}

//	func HashPassword(password string) string {
//		bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
//		return string(bytes)
//	}
func (uc *userUseCase) ListProducts(page, perPage int) ([]domain.ProductResponse, utils.MetaData, error) {

	users, metaData, err := uc.userRepo.ListProducts(page, perPage)
	if err != nil {
		return users, metaData, err
	}
	return users, metaData, nil
}

//Wishlist

func (uc *userUseCase) AddWishlist(user_id, product_id string) error {
	err := uc.userRepo.AddWishlist(user_id, product_id)
	if err != nil {
		return err
	}
	return nil
}
func (uc *userUseCase) ViewWishList(user_id string, page, perPage int) ([]domain.WishListResponse, utils.MetaData, error) {
	wishListResponse, metaData, err := uc.userRepo.ViewWishList(user_id, page, perPage)
	if err != nil {
		return wishListResponse, metaData, err
	}
	return wishListResponse, metaData, nil
}
func (uc *userUseCase) DeleteWishList(wishlist_id string) error {
	err := uc.userRepo.DeleteWishList(wishlist_id)
	if err != nil {
		return err
	}
	return nil
}

//Cart

func (uc *userUseCase) AddCart(user_id, product_id string) (error, string) {
	err, cart_id := uc.userRepo.AddCart(user_id, product_id)
	if err != nil {
		return err, cart_id
	}
	return nil, cart_id
}
func (uc *userUseCase) ViewCart(user_id string, page, perPage int) ([]domain.CartResponse, float32, utils.MetaData, error) {
	carts, grand_total, metaData, err := uc.userRepo.ViewCart(user_id, page, perPage)
	if err != nil {
		return carts, 0, metaData, err
	}
	return carts, grand_total, metaData, nil
}
func (uc *userUseCase) ListCoupon(user_id, product_id string) ([]domain.CouponResponse, error) {
	coupons, err := uc.userRepo.ListCoupon(user_id, product_id)
	if err != nil {
		return coupons, err
	}
	return coupons, nil
}
func (uc *userUseCase) ApplyCoupon(cart_id, order_id, coupon_id string) error {

	err := uc.userRepo.ApplyCoupon(cart_id, order_id, coupon_id)
	if err != nil {
		return err
	}
	return nil

}
func (uc *userUseCase) CancelCoupon(cart_id, order_id, coupon_id string) error {

	err := uc.userRepo.CancelCoupon(cart_id, order_id, coupon_id)
	if err != nil {
		return err
	}
	return nil

}
func (uc *userUseCase) DeleteCart(user_id, product_id string) error {
	err := uc.userRepo.DeleteCart(user_id, product_id)
	if err != nil {
		return err
	}
	return nil
}

//Shipping

func (uc *userUseCase) AddShippingDetails(user_id string, newAddress domain.ShippingDetails) error {
	err := uc.userRepo.AddShippingDetails(user_id, newAddress)
	if err != nil {
		return err
	}
	return nil
}
func (uc *userUseCase) ListShippingDetails(user_id string) ([]domain.ShippingDetailsResponse, error) {
	shippingDetails, err := uc.userRepo.ListShippingDetails(user_id)
	if err != nil {
		return shippingDetails, err
	}
	return shippingDetails, nil
}
func (uc *userUseCase) DeleteShippingDetails(user_id, address_id string) error {
	err := uc.userRepo.DeleteShippingDetails(user_id, address_id)
	if err != nil {
		return err
	}
	return nil
}

// Order
func (uc *userUseCase) CheckOut(cart_id, user_id, product_id, address_id string) (string, error) {
	if product_id == "" && cart_id == "" {
		return "", errors.New("Please select a product")
	}
	if address_id == "" {
		return "", errors.New("Please enter the shipping details")
	}
	id, err := uc.userRepo.CheckOut(cart_id, user_id, product_id, address_id)
	if err != nil {
		return "", err
	}
	return id, nil
}
func (uc *userUseCase) OrderSummery(order_id string) ([]domain.OrderSummery, error) {
	orderSummery, err := uc.userRepo.OrderSummery(order_id)
	if err != nil {
		return orderSummery, err
	}
	return orderSummery, err
}
func (uc *userUseCase) UpdateOrder(orders_id, product_id string, orderUpdates interface{}) error {
	_, err := uc.userRepo.OrderSummery(orders_id)
	if err != nil {
		return errors.New("No order found in checkout")
	}
	err = uc.userRepo.UpdateOrder(orders_id, product_id, orderUpdates)
	if err != nil {
		return err
	}
	return nil
}
