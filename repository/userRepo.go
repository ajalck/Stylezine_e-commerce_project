package repository

import (
	"ajalck/e_commerce/domain"
	repoInt "ajalck/e_commerce/repository/interface"
	"ajalck/e_commerce/utils"
	"errors"
	"fmt"

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
	var Products []domain.ProductResponse
	var totalRecords int64

	ur.DB.Model(&domain.Products{}).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return Products, metaData, err
	}

	result := ur.DB.Model(&domain.Products{}).Select("id", "item", "product_name", "discription", "product_image", "category_name", "brand_name", "size", "color", "unit_price", "stock", "rating").
		Joins("inner join categories on products.category_id=categories.category_id").
		Joins("inner join brands on products.brand_id=brands.brand_id").Offset(offset).Limit(perPage).Find(&Products)
	is := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if is == true {
		return Products, metaData, errors.New("Record not found")
	}
	return Products, metaData, nil
}
func (ur *UserRepo) ViewProduct(id int) (domain.ProductResponse, error) {

	result := ur.DB.Model(&domain.Products{}).Where("id", id).Where("status", "available").First(&domain.Products{})
	if is := errors.Is(result.Error, gorm.ErrRecordNotFound); is == true {
		fmt.Println("error is ", result.Error.Error())
		return domain.ProductResponse{}, result.Error
	}
	fmt.Println(domain.ProductResponse{})
	return domain.ProductResponse{}, nil
}

//Wish List

func (ur *UserRepo) AddWishlist(user_id, product_id int) error {

	wishlist := domain.WishList{
		User_ID:    user_id,
		Product_ID: product_id,
	}
	result := ur.DB.Where("id", product_id).First(&domain.Products{})
	if result.Error != nil {
		return errors.New("Product not found")
	}
	result = ur.DB.Where(&domain.WishList{User_ID: user_id, Product_ID: product_id}).First(&domain.WishList{})
	if is := errors.Is(result.Error, gorm.ErrRecordNotFound); is == false {
		return errors.New("Selected Item is already in your wishlist")
	}
	result = ur.DB.Select("user_id", "product_id").Create(&wishlist)
	if is := errors.Is(result.Error, gorm.ErrRegistered); is == true {
		return result.Error
	}
	return nil
}

func (ur *UserRepo) ViewWishList(user_id, page, perPage int) ([]domain.WishListResponse, utils.MetaData, error) {

	var favourites []domain.WishListResponse
	var totalRecords int64

	ur.DB.Model(&domain.WishList{}).Where("user_id", user_id).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return favourites, metaData, err
	}
	result := ur.DB.Model(&domain.Products{}).Select("id", "item", "product_name", "product_image", "size", "color", "status").
		Joins("right join wish_lists on products.id=wish_lists.product_id").Where("wish_lists.user_id", user_id).Offset(offset).Limit(perPage).Find(&favourites)

	if is := errors.Is(result.Error, gorm.ErrRecordNotFound); is == true {
		return favourites, metaData, result.Error
	}

	return favourites, metaData, nil
}
func (ur *UserRepo) DeleteWishList(user_id, product_id int) error {

	wishlist := domain.WishList{}
	result := ur.DB.Where(&domain.WishList{User_ID: user_id, Product_ID: product_id}).First(&wishlist)
	if result.Error != nil {
		return result.Error
	}
	ur.DB.Where(&domain.WishList{User_ID: user_id, Product_ID: product_id}).Delete(&wishlist)
	return nil
}

//Cart

func (ur *UserRepo) AddCart(user_id, product_id int) error {

	product := &domain.Products{}

	ur.DB.Table("products").Select("unit_price").Where("id", product_id).First(&product)
	unit_price := product.Unit_Price
	cart := &domain.Cart{
		User_ID:     user_id,
		Product_ID:  product_id,
		Count:       1,
		Total_Price: unit_price,
	}
	Cart, err := ur.CheckExistency(user_id, product_id)
	if err == nil {
		Cart.Count = Cart.Count + 1
		Cart.Total_Price = float32(Cart.Count) * unit_price
		ur.DB.Model(&cart).Where(&domain.Cart{User_ID: user_id, Product_ID: product_id}).Updates(&domain.Cart{Count: Cart.Count, Total_Price: Cart.Total_Price})
		return nil
	}

	result := ur.DB.Select("user_id", "product_id", "count", "total_price").Create(&cart)
	if is := errors.Is(result.Error, gorm.ErrRegistered); is == true {
		return result.Error
	}
	return nil
}

func (ur *UserRepo) CheckExistency(user_id, product_id int) (*domain.Cart, error) {

	cart := &domain.Cart{}
	result := ur.DB.Where(&domain.Cart{User_ID: user_id, Product_ID: product_id}).First(&cart)
	return cart, result.Error
}

func (ur *UserRepo) ViewCart(user_id, page, perPage int) ([]domain.CartResponse, utils.MetaData, error) {

	var cart []domain.CartResponse
	var totalRecords int64

	ur.DB.Model(&domain.Cart{}).Where("user_id", user_id).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return cart, metaData, err
	}
	result := ur.DB.Model(&domain.Products{}).Select("user_id", "product_id", "item", "product_name", "product_image", "size", "color", "count", "total_price", "status").
		Joins("right join carts on products.id=carts.product_id").Where("carts.user_id", user_id).Offset(offset).Limit(perPage).Find(&cart)

	if is := errors.Is(result.Error, gorm.ErrRecordNotFound); is == true {
		return cart, metaData, result.Error
	}
	fmt.Println(cart)
	return cart, metaData, nil
}

func (ur *UserRepo) DeleteCart(user_id, product_id int) error {
	cart := &domain.Cart{}
	Cart, err := ur.CheckExistency(user_id, product_id)
	if err == nil {
		unit_price := (Cart.Total_Price / float32(Cart.Count))
		if Cart.Count > 1 {
			Cart.Count = Cart.Count - 1
			Cart.Total_Price = unit_price * float32(Cart.Count)
			ur.DB.Model(&cart).Where(&domain.Cart{User_ID: user_id, Product_ID: product_id}).Updates(&domain.Cart{Count: Cart.Count, Total_Price: Cart.Total_Price})
			return nil
		}
		result := ur.DB.Where(&domain.Cart{User_ID: user_id, Product_ID: product_id}).Delete(&cart)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
	return err
}

//Shipping

func (ur *UserRepo) AddShippingDetails(user_id int, newAddress domain.ShippingDetails) error {

	address := &domain.ShippingDetails{}
	result := ur.DB.Where(&domain.ShippingDetails{Address: newAddress.Address, User_ID: uint(user_id)}).First(&address)
	if result.Error == nil {
		return errors.New("Entered input is already one of your shipping details")
	}
	result = ur.DB.Create(&domain.ShippingDetails{First_Name: newAddress.First_Name,
		Last_Name: newAddress.Last_Name,
		Email:     newAddress.Email,
		Phone:     newAddress.Phone,
		City:      newAddress.City,
		Street:    newAddress.Street,
		Address:   newAddress.Address,
		Pin_code:  newAddress.Pin_code,
		Land_Mark: newAddress.Land_Mark,
		User_ID:   uint(user_id)})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (ur *UserRepo) ListShippingDetails(user_id int) ([]domain.ShippingDetailsResponse, error) {
	ShippingAddRes := []domain.ShippingDetailsResponse{}
	var totalRecords int64
	result := ur.DB.Table("shipping_details").Where(&domain.ShippingDetails{User_ID: uint(user_id)}).
		Select("id", "user_id", "concat(first_name,' ',last_name)as name", "email", "phone", "city", "street", "address", "pin_code", "land_mark").
		Where("deleted_at IS NULL").Find(&ShippingAddRes).Count(&totalRecords)

	if result.Error != nil {
		return nil, result.Error
	}
	if totalRecords == 0 {
		return nil, errors.New("No records found")
	}
	return ShippingAddRes, nil
}
func (ur *UserRepo) DeleteShippingDetails(user_id, address_id int) error {
	shipping_dtl := &domain.ShippingDetails{}
	fmt.Println(user_id)
	result := ur.DB.Where("id", address_id).Delete(&shipping_dtl)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (ur *UserRepo) PlaceOrder(user_id, product_id, address_id, coupon_id int) error {
	product, err := ur.ViewProduct(product_id)
	if err != nil {
		return err
	}
	result := ur.DB.Where("id", address_id).Find(&domain.ShippingDetails{})
	if result.Error != nil {
		return result.Error
	}
	result=ur.DB.Where("id")
}
