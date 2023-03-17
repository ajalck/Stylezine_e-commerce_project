package repository

import (
	"ajalck/e_commerce/domain"
	repoInt "ajalck/e_commerce/repository/interface"
	"ajalck/e_commerce/utils"
	"errors"
	"fmt"
	"time"

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
func (ur *UserRepo) ViewProduct(id int) (domain.Products, error) {
	product := domain.Products{}
	result := ur.DB.Model(&domain.Products{}).Where("id", id).Where("status", "available").First(&product)
	if is := errors.Is(result.Error, gorm.ErrRecordNotFound); is == true {
		fmt.Println("error is ", result.Error.Error())
		return product, result.Error
	}
	return product, nil
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

// Cart
func (ur *UserRepo) AddCart(user_id, product_id int) (error, string) {

	product := &domain.Products{}

	ur.DB.Table("products").Select("unit_price").Where("id", product_id).First(&product)
	unit_price := product.Unit_Price

	//
	excart := &domain.Cart{}
	result := ur.DB.Where(&domain.Cart{User_ID: user_id}).First(&excart)
	var id string
	if result.Error == nil {
		id = excart.Cart_ID
	} else {
		id = utils.GenerateID()
	}
	//

	cart := &domain.Cart{
		Cart_ID:     id,
		User_ID:     user_id,
		Product_ID:  product_id,
		Coupon_id:   0,
		Quantity:    1,
		Unit_Price:  unit_price,
		Total_Price: unit_price,
	}
	Cart, err := ur.CheckExistency(user_id, product_id)
	if err == nil {
		Cart.Quantity = Cart.Quantity + 1
		Cart.Total_Price = float32(Cart.Quantity) * unit_price
		ur.DB.Model(&cart).Where(&domain.Cart{User_ID: user_id, Product_ID: product_id}).Updates(&domain.Cart{Quantity: Cart.Quantity, Total_Price: Cart.Total_Price})
		return nil, Cart.Cart_ID
	}
	result = ur.DB.Select("cart_id", "user_id", "product_id", "coupon_id", "quantity", "unit_price", "total_price").Create(&cart)
	if is := errors.Is(result.Error, gorm.ErrRegistered); is == true {
		return result.Error, cart.Cart_ID
	}
	return nil, cart.Cart_ID
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
func (ur *UserRepo) UpdateCoupon() error {
	coupons := []domain.Coupon{}
	ur.DB.Find(&coupons)
	for i := range coupons {
		if coupons[i].Expires_At.Compare(time.Now()) == -1 {
			result := ur.DB.Table("coupons").Where("id", coupons[i].ID).Update("coupon_status", "expired")
			if result.Error != nil {
				return result.Error
			}
		}
	}
	return nil
}
func (ur *UserRepo) DeleteCart(user_id, product_id int) error {
	cart := &domain.Cart{}
	Cart, err := ur.CheckExistency(user_id, product_id)
	if err == nil {
		unit_price := (Cart.Total_Price / float32(Cart.Quantity))
		if Cart.Quantity > 1 {
			Cart.Quantity = Cart.Quantity - 1
			Cart.Total_Price = unit_price * float32(Cart.Quantity)
			ur.DB.Model(&cart).Where(&domain.Cart{User_ID: user_id, Product_ID: product_id}).Updates(&domain.Cart{Quantity: Cart.Quantity, Total_Price: Cart.Total_Price})
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

//Coupon

func (ur *UserRepo) ListCoupon(user_id, product_id int) ([]domain.CouponResponse, error) {
	coupons := []domain.CouponResponse{}
	err := ur.UpdateCoupon()
	if err != nil {
		return coupons, err
	}
	result := ur.DB.Raw("SELECT id,coupon_code,discount_amount,user_id,product_id,min_cost,expires_at,coupon_status FROM coupons WHERE user_id IN ($1,0) AND product_id IN ($2,0) AND coupon_status=$3;", user_id, product_id, "active").Scan(&coupons)
	if result.Error != nil {
		return coupons, result.Error
	}
	return coupons, nil
}
func (ur *UserRepo) ValidateCoupon(user_id, product_id, coupon_id int) (bool, error) {
	user := domain.User{}
	ur.DB.Where("id", user_id).First(&user)
	product := domain.Products{}
	ur.DB.Where("id", product_id).First(&product)
	coupon := domain.Coupon{}
	err := ur.DB.Where("id", coupon_id).First(&coupon)
	if err.Error != nil {
		return false, errors.New("Coupon does not exists !")
	}
	if coupon.Coupon_Status == "active" {
		if coupon.Coupon_Code == "WELCOME200" {
			if user.Level == "bronze" {
				return true, nil
			}
			return false, errors.New("Coupon alredy used")
		} else {
			orders := []domain.Order{}
			result := ur.DB.Where("user_id", user_id).Find(&orders)
			if result.Error != nil {
				for i := range orders {
					if orders[i].Coupon_ID == uint(coupon_id) {
						return false, errors.New("Coupon alredy used")
					}
				}
			}
		}
		return true, nil
	}
	return false, errors.New("Coupon expired")
}
func (ur *UserRepo) ApplyCoupon(cart_id, order_id string, coupon_id int) error {
	if cart_id == "" {
		order := &domain.Order{}
		ur.DB.Where("id", order_id).First(&order)
		if order.Order_ID != "" {
			if order.Coupon_ID != 0 {
				return errors.New("one coupon is already applied")
			}
			valid, err := ur.ValidateCoupon(int(order.User_ID), int(order.Product_ID), coupon_id)
			if valid == true {
				result := ur.DB.Table("orders").Where("order_id", order_id).Update("coupon_id", coupon_id)
				if result.Error != nil {
					return result.Error
				}
			} else {
				return err
			}
		} else {
			return errors.New("No product in the checklist")
		}
	} else {
		cart := &domain.Cart{}
		ur.DB.Where("cart_id", cart_id).First(&cart)
		if cart.Cart_ID != "" {
			if cart.Coupon_id != 0 {
				return errors.New("one coupon is already applied")
			}
			valid, err := ur.ValidateCoupon(int(cart.User_ID), int(cart.Product_ID), coupon_id)
			if valid == true {
				result := ur.DB.Table("carts").Where("cart_id", cart_id).Update("coupon_id", coupon_id)
				if result.Error != nil {
					return result.Error
				}
			} else {
				return err
			}
		} else {
			return errors.New("No product in the cart")
		}
	}
	return nil
}
func (ur *UserRepo) CancelCoupon(cart_id, order_id string, coupon_id int) error {
	if cart_id == "" {
		order := &domain.Order{}
		ur.DB.Where("id", order_id).First(&order)
		if order.Order_ID != "" {
			if order.Coupon_ID == 0 {
				return errors.New("not found applied coupons !")
			}

			result := ur.DB.Model(&domain.Cart{}).Where("id", order_id).Delete("carts.coupon_id", coupon_id)
			if result.Error != nil {
				return result.Error
			}

		} else {
			return errors.New("No product in the checklist")
		}
	} else {
		cart := &domain.Cart{}
		ur.DB.Where("cart_id", cart_id).First(&cart)
		if cart.Cart_ID != "" {
			if cart.Coupon_id == 0 {
				return errors.New("not found applied coupons !")
			}

			result := ur.DB.Model(&domain.Cart{}).Where("cart_id", cart_id).Delete("carts.coupon_id", coupon_id)
			if result.Error != nil {
				return result.Error
			}

		} else {
			return errors.New("No product in the cart")
		}
	}
	return nil
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
func (ur *UserRepo) CheckOut(cart_id string, user_id, product_id, address_id int) (string, error) {
	result := ur.DB.Where("id", address_id).First(&domain.ShippingDetails{})
	if result.Error != nil {
		return "", result.Error
	}
	if cart_id != "" {
		cart := []domain.Cart{}
		result := ur.DB.Where("cart_id", cart_id).Find(&cart)
		if result.Error != nil {
			return "", result.Error
		}
		id := utils.GenerateID()
		var totalPrice float32 = 0
		var discount float32 = 0
		for i := range cart {
			coupon := &domain.Coupon{}
			if cart[i].Coupon_id != 0 {
				result := ur.DB.Where("id", cart[i].Coupon_id).First(&coupon)
				if result.Error != nil {
					errors.Join(result.Error)
				}
			}

			totalPrice = totalPrice + cart[i].Total_Price
			discount = discount + coupon.Discount_amount
			result = ur.DB.Create(&domain.Order{
				Order_ID:       id,
				User_ID:        uint(cart[i].User_ID),
				Product_ID:     uint(cart[i].Product_ID),
				Shipping_ID:    uint(address_id),
				Coupon_ID:      uint(cart[i].Coupon_id),
				Quantity:       cart[i].Quantity,
				Discount:       coupon.Discount_amount,
				TotalPrice:     totalPrice - coupon.Discount_amount,
				Grand_Total:    totalPrice - discount,
				GST:            ((totalPrice - discount) * 12) / 100,
				Final:          (totalPrice - discount) + (((totalPrice - discount) * 12) / 100),
				Order_Status:   "pending",
				Payment_Status: "pending",
			})
			if result.Error != nil {
				return "", result.Error
			}

		}
		return id, nil
	} else {
		product, err := ur.ViewProduct(product_id)
		if err != nil {
			return "", err
		}
		id := utils.GenerateID()
		result = ur.DB.Create(&domain.Order{
			Order_ID:       id,
			User_ID:        uint(user_id),
			Product_ID:     uint(product_id),
			Shipping_ID:    uint(address_id),
			Coupon_ID:      0,
			Quantity:       1,
			Discount:       0,
			TotalPrice:     product.Unit_Price,
			Order_Status:   "pending",
			Payment_Status: "pending",
		})
		if result.Error != nil {
			return "", result.Error
		}
		return id, nil
	}
}
func (ur *UserRepo) OrderSummery(order_id string) ([]domain.OrderSummery, error) {
	orderSummery := []domain.OrderSummery{}
	// var totalRecords int64
	// page := 1
	// perPage := 1
	// ur.DB.Model(&domain.Order{}).Where("user_id", user_id).Count(&totalRecords)
	// metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	// if err != nil {
	// 	return orderSummery, metaData, err
	// }
	// results := ur.DB.Model(&domain.Order{}).Select("concat(shipping_details.first_name,' ',shipping_details.last_name) as shipping_name", "address as shipping_address", "product_name", "discription",
	// 	"product_image", "quantity", "sum(discount)as discount", "sum(total_price)as grand_total", "mode_of_payment", "order_status", "payment_status").
	// 	Joins("left join shipping_details on shipping_details.id=orders.shipping_id").
	// 	Joins("left join products on products.id=orders.product_id").
	// 	Where("orders.user_id", user_id).Offset(offset).Limit(perPage).Find(&orderSummery)

	results := ur.DB.Model(&domain.Order{}).Where(&domain.Order{Order_ID: order_id, Order_Status: "pending"}).Find(&orderSummery)
	if results.Error != nil {
		return orderSummery, results.Error
	}
	return orderSummery, results.Error
}
func (ur *UserRepo) UpdateOrder(order_id, product_id string, orderUpdates interface{}) error {
	// order := domain.Order{}
	// result := ur.DB.Where(domain.Order{Order_ID: order_id, Product_ID: product_id}).First(&order)
	// if result.Error != nil {
	// 	return result.Error
	// }
	return nil
}
