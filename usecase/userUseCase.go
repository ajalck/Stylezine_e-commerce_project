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
	newUser.Verification = false
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
