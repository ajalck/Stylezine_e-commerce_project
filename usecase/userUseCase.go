package usecase

import (
	"ajalck/e_commerce/domain"
	repoInt "ajalck/e_commerce/repository/interface"
	services "ajalck/e_commerce/usecase/interface"
	"ajalck/e_commerce/utils"
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userUseCase struct {
	userRepo repoInt.UserRepository
}

func NewUserUseCase(repo repoInt.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}
func (uc *userUseCase) CreateUser(c *gin.Context, newUser domain.User) error {

	newUser.User_Role = "user"

	if _, err := uc.userRepo.FindUser(c, newUser.Email, newUser.User_Role); err == nil {

		c.JSON(http.StatusBadGateway, gin.H{"message": "User Already exists"})
		err = errors.New("user already exists")
		return err

	}

	//hashing password

	newUser.Password = HashPassword(newUser.Password)
	newUser.Status = "active"
	newUser.Level = "bronze"
	err := uc.userRepo.CreateUser(c, newUser)
	if err != nil {
		c.JSON(400, gin.H{"error": "couldn't create a user"})
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

func (uc *userUseCase) AddWishlist(user_id, product_id int) error {
	err := uc.userRepo.AddWishlist(user_id, product_id)
	if err != nil {
		return err
	}
	return nil
}
func (uc *userUseCase) ViewWishList(user_id, page, perPage int) ([]domain.WishListResponse, utils.MetaData, error) {
	wishListResponse, metaData, err := uc.userRepo.ViewWishList(user_id, page, perPage)
	if err != nil {
		return wishListResponse, metaData, err
	}
	return wishListResponse, metaData, nil
}
func (uc *userUseCase) DeleteWishList(user_id, product_id int) error {
	err := uc.userRepo.DeleteWishList(user_id, product_id)
	if err != nil {
		return err
	}
	return nil
}

//Cart

func (uc *userUseCase) AddCart(user_id, product_id int) error {
	err := uc.userRepo.AddCart(user_id, product_id)
	if err != nil {
		return err
	}
	return nil
}
func (uc *userUseCase) ViewCart(user_id, page, perPage int) ([]domain.CartResponse, utils.MetaData, error) {
	carts, metaData, err := uc.userRepo.ViewCart(user_id, page, perPage)
	if err != nil {
		return carts, metaData, err
	}
	return carts, metaData, nil
}
func (uc *userUseCase) ListCoupon(user_id, product_id int) ([]domain.CouponResponse, error) {
	coupons, err := uc.userRepo.ListCoupon(user_id, product_id)
	if err != nil {
		return coupons, err
	}
	return coupons, nil
}
func (uc *userUseCase) ApplyCoupon(cart_id, order_id, coupon_id int) error {

	err := uc.userRepo.ApplyCoupon(cart_id, order_id, coupon_id)
	if err != nil {
		return err
	}
	return nil

}
func (uc *userUseCase) CancelCoupon(cart_id, order_id, coupon_id int) error {

	err := uc.userRepo.CancelCoupon(cart_id, order_id, coupon_id)
	if err != nil {
		return err
	}
	return nil

}
func (uc *userUseCase) DeleteCart(user_id, product_id int) error {
	err := uc.userRepo.DeleteCart(user_id, product_id)
	if err != nil {
		return err
	}
	return nil
}

//Shipping

func (uc *userUseCase) AddShippingDetails(user_id int, newAddress domain.ShippingDetails) error {
	err := uc.userRepo.AddShippingDetails(user_id, newAddress)
	if err != nil {
		return err
	}
	return nil
}
func (uc *userUseCase) ListShippingDetails(user_id int) ([]domain.ShippingDetailsResponse, error) {
	shippingDetails, err := uc.userRepo.ListShippingDetails(user_id)
	if err != nil {
		return shippingDetails, err
	}
	return shippingDetails, nil
}
func (uc *userUseCase) DeleteShippingDetails(user_id, address_id int) error {
	err := uc.userRepo.DeleteShippingDetails(user_id, address_id)
	if err != nil {
		return err
	}
	return nil
}

// Order
func (uc *userUseCase) CheckOut(user_id, cart_id, product_id, address_id int) error {
	if product_id == 0 && cart_id == 0 {
		return errors.New("Please select a product")
	}
	if address_id == 0 {
		return errors.New("Please enter the shipping details")
	}
	err := uc.userRepo.CheckOut(user_id, cart_id, product_id, address_id)
	if err != nil {
		return err
	}
	return nil
}
