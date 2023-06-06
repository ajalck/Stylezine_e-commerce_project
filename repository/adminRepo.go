package repository

import (
	"ajalck/e_commerce/domain"
	repoInt "ajalck/e_commerce/repository/interface"
	"ajalck/e_commerce/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AdminRepo struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) repoInt.AdminRepository {
	return &AdminRepo{DB: DB}
}
func (ar *AdminRepo) CreateAdmin(newAdmin domain.User) error {

	err := ar.DB.Create(&newAdmin)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (ar *AdminRepo) FindAdmin(email string, userRole string) (domain.User, error) {

	var admin domain.User
	// user := ar.DB.First(&admin, "Email=?", email)

	// user := ar.DB.Where("Email = ? AND UserRole = ?", email, userRole).First(&admin)

	user := ar.DB.Where(&domain.User{Email: email, User_Role: userRole}).First(&admin)

	if user.Error != nil {
		return admin, errors.New("could'nt find admin")
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

	result := ar.DB.Model(&domain.User{}).Select("user_id", "first_name", "last_name", "email", "gender", "phone", "status", "user_role").Where("user_role", "user").Offset(offset).Limit(perPage).Find(&Users)
	is := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if is == true {
		return Users, metaData, errors.New("Record not found")
	}
	return Users, metaData, nil
}
func (ar *AdminRepo) ViewUser(id string) (domain.UserResponse, error) {

	User := domain.UserResponse{}
	ar.DB.Table("users").Select("user_id", "first_name", "last_name", "email", "gender", "phone", "status", "user_role").Where("user_id=?", id).Where("user_role", "user").Find(&User)
	if User.User_ID == "" {
		err := errors.New("no user found")
		return User, err
	}
	return User, nil
}
func (ar *AdminRepo) BlockUser(id string) {
	var user domain.User
	ar.DB.Raw("UPDATE users SET status=$1 WHERE user_id=$2;", "blocked", id).Scan(&user)

}
func (ar *AdminRepo) UnblockUser(id string) {
	var user domain.User
	ar.DB.Raw("UPDATE users SET status=$1 WHERE user_id=$2;", "active", id).Scan(&user)

}
func (ar *AdminRepo) ListBlockedUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error) {

	var Users []domain.UserResponse
	var totalRecords int64

	ar.DB.Model(&domain.User{}).Where(&domain.User{User_Role: "user", Status: "blocked"}).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return Users, metaData, err
	}

	result := ar.DB.Model(&domain.User{}).Select("user_id", "first_name", "last_name", "email", "gender", "phone", "status", "user_role").Where(&domain.User{User_Role: "user", Status: "blocked"}).Offset(offset).Limit(perPage).Find(&Users)
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

	result := ar.DB.Model(&domain.User{}).Select("user_id", "first_name", "last_name", "email", "gender", "phone", "status", "user_role").Where(&domain.User{User_Role: "user", Status: "active"}).Offset(offset).Limit(perPage).Find(&Users)
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

	result := ar.DB.Where("category_name", category.Category_name).First(&domain.Category{})
	if result.Error == nil {
		return errors.New("A Category already exists on the same name")
	}
	var category_id int
	if err := ar.DB.Raw("INSERT INTO categories (category_name) VALUES ($1) RETURNING category_id;", category.Category_name).Scan(&category_id).Error; err != nil {
		return err
	}
	return nil

}
func (ar *AdminRepo) ListCategory() ([]domain.Category, error) {
	var categories []domain.Category
	dbResult := ar.DB.Find(&categories)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return categories, dbResult.Error
	}
	return categories, nil
}
func (ar *AdminRepo) EditCategory(category domain.Category) error {

	err := ar.DB.Raw("UPDATE categories SET category_name=$1 WHERE category_id=$2;", category.Category_name, category.Category_ID).Scan(&category).Error
	if err != nil {
		return err
	}

	return nil

}
func (ar *AdminRepo) DeleteCategory(category domain.Category) error {

	err := ar.DB.Where(&domain.Category{Category_ID: category.Category_ID}).Delete(&category)
	if err.Error != nil {
		return err.Error
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

	result := ar.DB.Where("brand_name", brand.Brand_Name).First(&domain.Brand{})
	if result.Error == nil {
		return errors.New("A Brand already exists on the same name")
	}
	if err := ar.DB.Raw("INSERT INTO brands (brand_name,brand_discription) VALUES ($1,$2);", brand.Brand_Name, brand.Brand_Discription).Scan(&brand).Error; err != nil {
		return err
	}
	return nil
}
func (ar *AdminRepo) ListBrands() ([]domain.Brand, error) {

	var brands []domain.Brand
	dbResult := ar.DB.Find(&brands)
	if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		return brands, dbResult.Error
	}
	return brands, nil
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

func (ar *AdminRepo) AddProducts(products domain.Products) (string, error) {
	product := domain.Products{}
	results := ar.DB.Where(&domain.Products{
		Product_Name: products.Product_Name,
		Category_id:  products.Category_id,
		Brand_id:     products.Brand_id,
		Size:         products.Size,
		Color:        products.Color}).First(&product)
	if results.Error != nil {
		if err := ar.DB.Raw("INSERT INTO products (product_code,item,product_name,discription,product_image,category_id,brand_id,size,color,unit_price,stock,status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)RETURNING product_code ;",
			utils.GenerateID(), products.Item, products.Product_Name, products.Discription, products.Product_Image, products.Category_id, products.Brand_id, products.Size, products.Color, products.Unit_Price, products.Stock, products.Status).Scan(&products).Error; err != nil {
			return "", err
		}
		return products.Product_Code, nil
	}
	err := ar.DB.Model(&domain.Products{}).Where("product_code", product.Product_Code).Update("stock", (*product.Stock + 1))
	if err.Error != nil {
		return "", err.Error
	}
	return product.Product_Code, nil
}
func (ar *AdminRepo) ListProducts(page, perPage int) ([]domain.ProductResponse, utils.MetaData, error) {
	var Products []domain.ProductResponse
	var totalRecords int64

	ar.DB.Model(&domain.Products{}).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return Products, metaData, err
	}

	result := ar.DB.Model(&domain.Products{}).Select("product_code", "item", "product_name", "discription", "product_image", "category_name", "brand_name", "size", "color", "unit_price", "stock", "rating", "status").
		Joins("inner join categories on products.category_id=categories.category_id").
		Joins("inner join brands on products.brand_id=brands.brand_id").Offset(offset).Limit(perPage).Find(&Products)
	is := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if is == true {
		return Products, metaData, errors.New("Record not found")
	}
	return Products, metaData, nil

}
func (ar *AdminRepo) EditProducts(product domain.Products) error {
	P := &domain.Products{}
	results := ar.DB.Where("product_code", product.Product_Code).First(&P)
	if results.Error != nil {
		return results.Error
	}

	results = ar.DB.Model(&domain.Products{}).Where("product_code", product.Product_Code).
		Updates(map[string]interface{}{
			"product_image": gorm.Expr("COALESCE(?, products.product_name)", product.Product_Image),
			"size":          gorm.Expr("COALESCE(?, products.size)", product.Size),
			"color":         gorm.Expr("COALESCE(?, products.color)", product.Color),
			"unit_price":    gorm.Expr("COALESCE(?, products.unit_price)", product.Unit_Price),
			"stock":         gorm.Expr("COALESCE(?, products.stock)", product.Stock),
		})
	if results.Error != nil {
		return results.Error
	}
	return nil

}
func (ar *AdminRepo) DeleteProducts(product_id string) error {

	if err := ar.DB.Raw("DELETE FROM products WHERE product_code=$1;", product_id).Scan(&domain.Products{}).Error; err != nil {
		return err
	}
	return nil
}
func (ar *AdminRepo) AddCoupon(coupon domain.Coupon) error {
	if result := ar.DB.Where(&domain.Coupon{Coupon_Code: coupon.Coupon_Code, Coupon_Status: "active"}).First(&domain.Coupon{}); result.Error == nil {
		return errors.New("Coupon already exists")
	}
	expiresAt := coupon.Expires_At
	duration := time.Until(expiresAt)
	fmt.Println(duration)
	if duration < 0 {
		coupon.Coupon_Status = "expired"
	} else {
		coupon.Coupon_Status = "active"
	}
	result := ar.DB.Create(&domain.Coupon{
		Coupon_Code:     coupon.Coupon_Code,
		Discount_amount: coupon.Discount_amount,
		User_ID:         coupon.User_ID,
		Product_ID:      coupon.Product_ID,
		Min_Cost:        coupon.Min_Cost,
		Expires_At:      coupon.Expires_At,
		Coupon_Status:   coupon.Coupon_Status,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil

}
func (ar *AdminRepo) ListCoupon(page, perPage int) ([]domain.CouponResponse, utils.MetaData, error) {
	var coupons []domain.CouponResponse
	var totalRecords int64
	coupon := domain.Coupon{}
	ar.DB.Find(&coupon).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))
	if err != nil {
		return coupons, metaData, err
	}
	results := ar.DB.Model(&coupon).Select("coupon_code", "discount_amount", "user_id", "product_id", "min_cost", "expires_at", "coupon_status").Offset(offset).Limit(perPage).Find(&coupons)
	if results.Error != nil {
		return coupons, metaData, results.Error
	}
	return coupons, metaData, nil
}
func (ar *AdminRepo) DeleteCoupon(coupon_id string) error {
	coupon := &domain.Coupon{}
	result := ar.DB.Where("coupon_code", coupon_id).First(&coupon)
	if result.Error == nil {
		result := ar.DB.Where("coupon_code", coupon_id).Unscoped().Delete(&coupon)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		return result.Error
	}
}

//Sales

func (ar *AdminRepo) SalesReport(page, perPage int) (interface{}, utils.MetaData, error) {
	var totalRecords int64

	ar.DB.Model(&domain.OrderReport{}).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return nil, metaData, err
	}
	type Sales_Report struct {
		Order_ID       string
		User_ID        string
		Product_ID     string
		Quantity       int
		TotalPrice     float32
		Payment_ID     string
		Order_Status   string
		Payment_Status string
	}
	sales_Report := []Sales_Report{}
	result := ar.DB.Model(&domain.OrderReport{}).Select("orders.order_id", "orders.user_id", "order_reports.product_id", "order_reports.quantity", "order_reports.total_price",
		"payment_id", "order_reports.order_status", "payment_status").Joins("right join orders on orders.order_id=order_reports.order_id").Offset(offset).Limit(perPage).Find(&sales_Report)
	if result.Error != nil {
		return nil, metaData, result.Error
	}
	return sales_Report, metaData, nil
}
