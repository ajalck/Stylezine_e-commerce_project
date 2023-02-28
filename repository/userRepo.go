package repository

import (
	"ajalck/e_commerce/domain"
	repoInt "ajalck/e_commerce/repository/interface"
	"ajalck/e_commerce/utils"
	"errors"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) repoInt.UserRepository {
	return &UserRepo{DB: DB}
}
func (ur *UserRepo) CreateUser(c *gin.Context, newUser domain.User) error {

	user := ur.DB.Create(&newUser)
	if user.Error != nil {
		return errors.New("couldn't create a new user")
	}
	return nil

}
func (ur *UserRepo) FindUser(c *gin.Context, email string, userRole string) (domain.User, error) {

	var users domain.User

	// user := ur.DB.First(&users, "Email=?", email)

	// user := ur.DB.Where("Email = ? AND UserRole = ?", email, userRole).Find(&users)

	user := ur.DB.Where(&domain.User{Email: email, User_Role: userRole}).First(&users)

	if user.Error != nil {
		return users, errors.New("could'nt find user")
	}
	return users, nil
}

func (ur *UserRepo) ListProducts(page, perPage int) ([]domain.ProductResponse, utils.MetaData, error) {
	var Users []domain.ProductResponse
	var totalRecords int64

	ur.DB.Model(&domain.Products{}).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return Users, metaData, err
	}

	result := ur.DB.Model(&domain.Products{}).Select("id", "item", "product_name", "discription", "product_image", "category_name", "brand_name", "size", "color", "unit_price", "stock", "rating").
		Joins("inner join categories on products.category_id=categories.category_id").
		Joins("inner join brands on products.brand_id=brands.brand_id").Offset(offset).Limit(perPage).Find(&Users)
	is := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if is == true {
		return Users, metaData, errors.New("Record not found")
	}
	return Users, metaData, nil
}
