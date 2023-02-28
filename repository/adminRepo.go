package repository

import (
	"ajalck/e_commerce/domain"
	repoInt "ajalck/e_commerce/repository/interface"
	"ajalck/e_commerce/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminRepo struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) repoInt.AdminRepository {
	return &AdminRepo{DB: DB}
}
func (ar *AdminRepo) CreateAdmin(c *gin.Context, newAdmin domain.User) error {

	err := ar.DB.Create(&newAdmin).Error
	return err
}
func (ar *AdminRepo) FindAdmin(c *gin.Context, email string, userRole string) (domain.User, error) {

	var admin domain.User
	// user := ar.DB.First(&admin, "Email=?", email)

	// user := ar.DB.Where("Email = ? AND UserRole = ?", email, userRole).First(&admin)

	user := ar.DB.Where(&domain.User{Email: email, User_Role: userRole}).First(&admin)

	if user.Error != nil {
		return admin, errors.New("could'nt find user")
	}
	return admin, nil
}

// User Management

func (ar *AdminRepo) ListUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error) {

	var Users []domain.UserResponse
	var totalRecords int64

	ar.DB.Model(&domain.User{}).Where("user_role", "user").Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return Users, metaData, err
	}

	result := ar.DB.Model(&domain.User{}).Select("id", "first_name", "last_name", "email", "gender", "phone", "status", "user_role").Where("user_role", "user").Offset(offset).Limit(perPage).Find(&Users)
	is := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if is == true {
		return Users, metaData, errors.New("Record not found")
	}
	return Users, metaData, nil
}
func (ar *AdminRepo) ViewUser(id int) (domain.UserResponse, error) {

	User := domain.UserResponse{}
	ar.DB.Raw("SELECT id,first_name,last_name,email,gender,phone,status,user_role FROM users WHERE id=?,user_role=?;", id, "user").Scan(&User)
	if User.ID == 0 {
		err := errors.New("no user found")
		return User, err
	}
	return User, nil
}
func (ar *AdminRepo) BlockUser(id int) {
	var user domain.User
	ar.DB.Raw("UPDATE users SET status=$1 WHERE id=$2;", "blocked", id).Scan(&user)
}
func (ar *AdminRepo) UnblockUser(id int) {
	var user domain.User
	ar.DB.Raw("UPDATE users SET status=$1 WHERE id=$2;", "active", id).Scan(&user)
}
func (ar *AdminRepo) ListBlockedUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error) {

	var Users []domain.UserResponse
	var totalRecords int64

	ar.DB.Model(&domain.User{}).Where(&domain.User{User_Role: "user", Status: "blocked"}).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return Users, metaData, err
	}

	result := ar.DB.Model(&domain.User{}).Select("id", "first_name", "last_name", "email", "gender", "phone", "status", "user_role").Where(&domain.User{User_Role: "user", Status: "blocked"}).Offset(offset).Limit(perPage).Find(&Users)
	is := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if is == true {
		return Users, metaData, errors.New("Record not found")
	}
	return Users, metaData, nil
}
func (ar *AdminRepo) ListActiveUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error) {

	var Users []domain.UserResponse
	var totalRecords int64

	ar.DB.Model(&domain.User{}).Where(&domain.User{User_Role: "user", Status: "active"}).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return Users, metaData, err
	}

	result := ar.DB.Model(&domain.User{}).Select("id", "first_name", "last_name", "email", "gender", "phone", "status", "user_role").Where(&domain.User{User_Role: "user", Status: "active"}).Offset(offset).Limit(perPage).Find(&Users)
	is := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if is == true {
		return Users, metaData, errors.New("Record not found")
	}
	return Users, metaData, nil
}

// Category Management

func (ar *AdminRepo) ViewCategory(category domain.Category) (domain.Category, error) {

	var categories domain.Category
	dbResult := ar.DB.Where("Category_ID=?", category.Category_ID).First(&categories)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return categories, dbResult.Error
	}
	return categories, nil
}

func (ar *AdminRepo) AddCategory(category domain.Category) error {

	var category_id int
	if err := ar.DB.Raw("INSERT INTO categories (category_name) VALUES ($1) RETURNING category_id;", category.Category_name).Scan(&category_id).Error; err != nil {
		return err
	}
	return nil

}
func (ar *AdminRepo) EditCategory(category domain.Category) error {

	err := ar.DB.Raw("UPDATE categories SET category_name=$1 WHERE category_id=$2;", category.Category_name, category.Category_ID).Scan(&category).Error
	if err != nil {
		return err
	}

	return nil

}
func (ar *AdminRepo) DeleteCategory(category domain.Category) error {

	if err := ar.DB.Raw("DELETE FROM categories WHERE category_id=$1 OR category_name=$2;", category.Category_ID, category.Category_name).Scan(&category).Error; err != nil {
		return err
	}
	return nil
}

// Brand Management

func (ar *AdminRepo) ViewBrand(brand_id uint) (domain.Brand, error) {

	var brands domain.Brand
	dbResult := ar.DB.Where("Brand_ID=?", brand_id).First(&brands)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return brands, dbResult.Error
	}
	return brands, nil
}

func (ar *AdminRepo) AddBrand(brand domain.Brand) error {

	if err := ar.DB.Raw("INSERT INTO brands (brand_name,brand_discription) VALUES ($1,$2);", brand.Brand_Name, brand.Brand_Discription).Scan(&brand).Error; err != nil {
		return err
	}
	return nil
}
func (ar *AdminRepo) EditBrand(brand domain.Brand) error {

	if err := ar.DB.Raw("UPDATE brands SET (brand_name,brand_discription)=($1,$2) WHERE brand_id=$3;", brand.Brand_Name, brand.Brand_Discription, brand.Brand_ID).Scan(&brand).Error; err != nil {
		return err
	}
	return nil

}
func (ar *AdminRepo) DeleteBrand(brand domain.Brand) error {

	if err := ar.DB.Raw("DELETE FROM brands WHERE brand_id =$1;", (brand.Brand_ID)).Scan(&brand).Error; err != nil {
		return err
	}
	return nil
}

// Product Management

func (ar *AdminRepo) AddProducts(products domain.Products) error {

	if err := ar.DB.Raw("INSERT INTO products (item,product_name,discription,product_image,category_id,brand_id,size,color,unit_price,stock) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)RETURNING id ;",
		products.Item, products.Product_Name, products.Discription, products.Product_Image, products.Category_id, products.Brand_id, products.Size, products.Color, products.Unit_Price, products.Stock).Scan(&products).Error; err != nil {
		return err
	}
	return nil
}
func (ar *AdminRepo) EditProducts(product domain.Products) error {

	if err := ar.DB.Raw("UPDATE products SET (product_image,UPPER(size),UPPER(color),unit_price,stock)=($1,$2,$3,$4,$5) WHERE brand_id=$6;",
		product.Product_Image, product.Size, product.Color, product.Unit_Price, product.Stock, product.ID).Scan(&product).Error; err != nil {
		return err
	}
	return nil

}
func (ar *AdminRepo) DeleteProducts(product domain.Products) error {

	if err := ar.DB.Raw("DELETE FROM products WHERE id=$1;", product.ID).Scan(&product).Error; err != nil {
		return err
	}
	return nil
}
