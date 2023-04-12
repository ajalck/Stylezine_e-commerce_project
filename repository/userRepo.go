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

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) repoInt.UserRepository {
	return &UserRepo{DB: DB}
}
func (ur *UserRepo) CreateUser(newUser domain.User) error {

	user := ur.DB.Create(&newUser)
	if user.Error != nil {
		return user.Error
	}
	return nil

}
func (ur *UserRepo) FindUser(email string, userRole string) (domain.User, error) {

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

	result := ur.DB.Model(&domain.Products{}).Select("product_code", "item", "product_name", "discription", "product_image", "category_name", "brand_name", "size", "color", "unit_price", "stock", "rating", "status").
		Joins("inner join categories on products.category_id=categories.category_id").
		Joins("inner join brands on products.brand_id=brands.brand_id").Offset(offset).Limit(perPage).Find(&Products)
	is := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if is == true {
		return Products, metaData, errors.New("Record not found")
	}
	return Products, metaData, nil
}
func (ur *UserRepo) ViewProduct(id string) (domain.Products, error) {
	product := domain.Products{}
	result := ur.DB.Model(&domain.Products{}).Where("product_code", id).Where("status", "available").First(&product)
	if is := errors.Is(result.Error, gorm.ErrRecordNotFound); is == true {
		fmt.Println("error is ", result.Error.Error())
		return product, result.Error
	}
	return product, nil
}

//Wish List

func (ur *UserRepo) AddWishlist(user_id, product_id string) error {

	wishlist := domain.WishList{
		Wishlist_ID: utils.GenerateID(),
		User_ID:     user_id,
		Product_ID:  product_id,
	}
	result := ur.DB.Where("product_code", product_id).First(&domain.Products{})
	if result.Error != nil {
		return errors.New("Product not found")
	}
	result = ur.DB.Where(&domain.WishList{User_ID: user_id, Product_ID: product_id}).First(&domain.WishList{})
	if is := errors.Is(result.Error, gorm.ErrRecordNotFound); is == false {
		return errors.New("Selected Item is already in your wishlist")
	}
	result = ur.DB.Select("wishlist_id", "user_id", "product_id").Create(&wishlist)
	if is := errors.Is(result.Error, gorm.ErrRegistered); is == true {
		return result.Error
	}
	return nil
}

func (ur *UserRepo) ViewWishList(user_id string, page, perPage int) ([]domain.WishListResponse, utils.MetaData, error) {

	var favourites []domain.WishListResponse
	var totalRecords int64
	ur.DB.Model(&domain.WishList{}).Where("user_id", user_id).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return favourites, metaData, err
	}
	result := ur.DB.Model(&domain.WishList{}).Select("wishlist_id", "user_id", "product_code", "item", "product_name", "product_image", "size", "color", "status").
		Joins("right join products on products.product_code=wish_lists.product_id").Where("wish_lists.user_id", user_id).Offset(offset).Limit(perPage).Find(&favourites)
	if is := errors.Is(result.Error, gorm.ErrRecordNotFound); is == true {
		return favourites, metaData, result.Error
	}
	return favourites, metaData, nil
}
func (ur *UserRepo) DeleteWishList(wishlist_id string) error {

	wishlist := domain.WishList{}
	result := ur.DB.Where(&domain.WishList{Wishlist_ID: wishlist_id}).First(&wishlist)
	if result.Error != nil {
		return result.Error
	}
	ur.DB.Where(&domain.WishList{Wishlist_ID: wishlist_id}).Delete(&wishlist)
	return nil
}

// Cart
func (ur *UserRepo) AddCart(user_id, product_id string) (error, string) {

	product := &domain.Products{}

	ur.DB.Table("products").Select("unit_price").Where("product_code", product_id).First(&product)
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
		Quantity:    1,
		Unit_Price:  *unit_price,
		Total_Price: *unit_price,
	}
	Cart, err := ur.CheckExistency(user_id, product_id)
	if err == nil {
		Cart.Quantity = Cart.Quantity + 1
		Cart.Total_Price = float32(Cart.Quantity) * *unit_price
		ur.DB.Model(&cart).Where(&domain.Cart{User_ID: user_id, Product_ID: product_id}).Updates(&domain.Cart{Quantity: Cart.Quantity, Total_Price: Cart.Total_Price})
		ur.Final_Cart(cart.Cart_ID, "not applied")
		return nil, Cart.Cart_ID
	}
	result = ur.DB.Select("cart_id", "user_id", "product_id", "quantity", "unit_price", "total_price").Create(&cart)
	if is := errors.Is(result.Error, gorm.ErrRegistered); is == true {
		return result.Error, cart.Cart_ID
	}
	ur.Final_Cart(cart.Cart_ID, "not applied")
	return nil, cart.Cart_ID
}

func (ur *UserRepo) CheckExistency(user_id, product_id string) (*domain.Cart, error) {

	cart := &domain.Cart{}
	result := ur.DB.Where(&domain.Cart{User_ID: user_id, Product_ID: product_id}).First(&cart)
	return cart, result.Error
}

func (ur *UserRepo) ViewCart(user_id string, page, perPage int) ([]domain.CartResponse, float32, utils.MetaData, error) {

	var cart []domain.CartResponse
	var totalRecords int64
	var grand_total float32 = 0
	ur.DB.Model(&domain.Cart{}).Where("user_id", user_id).Count(&totalRecords)
	metaData, offset, err := utils.ComputeMetaData(page, perPage, int(totalRecords))

	if err != nil {
		return cart, 0, metaData, err
	}
	result := ur.DB.Model(&domain.Cart{}).Select("carts.cart_id", "user_id", "product_id", "coupon_id", "item", "product_name", "product_image", "size", "color", "quantity", "carts.total_price", "status").
		Joins("right join products on products.product_code=carts.product_id").
		Joins("right join final_carts on carts.cart_id=final_carts.cart_id").
		Where("carts.user_id", user_id).Offset(offset).Limit(perPage).Find(&cart)

	if is := errors.Is(result.Error, gorm.ErrRecordNotFound); is == true {
		return cart, 0, metaData, result.Error
	}
	for i := range cart {
		grand_total = grand_total + cart[i].TotalPrice
	}
	return cart, grand_total, metaData, nil
}
func (ur *UserRepo) DeleteCart(user_id, product_id string) error {
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

// Coupon

func (ur *UserRepo) UpdateCoupon() error {
	coupons := []domain.Coupon{}
	ur.DB.Find(&coupons)
	for i := range coupons {
		result := coupons[i].Expires_At.Compare(time.Now())
		if result == -1 {
			result := ur.DB.Table("coupons").Where("id", coupons[i].ID).Update("coupon_status", "expired")
			if result.Error != nil {
				return result.Error
			}
		}
	}
	return nil
}

func (ur *UserRepo) ListCoupon(user_id, product_id string) ([]domain.CouponResponse, error) {
	coupons := []domain.CouponResponse{}
	err := ur.UpdateCoupon()
	if err != nil {
		return coupons, err
	}
	// sql := `select (coupon_code,discount_amount,user_id,product_id,min_cost,expires_at,coupon_status) from coupons where (user_id in($1) or user_id is null) and (product_id in($2)or product_id is null) and coupon_status ='active';`
	result := ur.DB.Model(&domain.Coupon{}).Where("(user_id = ? OR user_id IS NULL) AND (product_id = ? OR product_id IS NULL) AND coupon_status = ?", user_id, product_id, "active").Find(&coupons)
	if result.Error != nil {
		return coupons, result.Error
	}
	return coupons, nil
}
func (ur *UserRepo) ValidateCoupon(user_id, product_id, cart_id, coupon_id string) (bool, error) {
	user := domain.User{}
	ur.DB.Where("user_id", user_id).First(&user)
	product := domain.Products{}
	ur.DB.Where("product_code", product_id).First(&product)
	coupon := domain.Coupon{}
	err := ur.DB.Where("coupon_code", coupon_id).First(&coupon)
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
			result := ur.DB.Where("user_id=? AND order_status=? AND coupon_id=?", user_id, "success", coupon_id).First(&orders)
			if result.Error == nil {
				return false, errors.New("Coupon alredy used")
			}
			if coupon.Product_ID != nil {
				if product_id != "" {
					if coupon.Product_ID == &product_id {
						return true, nil
					}
				}
				return false, errors.New("This coupon is applicable for selected products only !")
			}
		}
		return true, nil
	}
	return false, errors.New("Coupon expired")
}
func (ur *UserRepo) ApplyCoupon(cart_id, order_id, coupon_id string) error {
	if cart_id == "" {
		order := &domain.Order{}
		orderreport := &domain.OrderReport{}
		ur.DB.Where("order_id", order_id).First(&order)
		ur.DB.Where("order_id", order_id).First(&orderreport)
		if order.Order_ID != "" {
			if order.Coupon_ID != "not applied" {
				return errors.New("one coupon is already applied")
			}
			valid, err := ur.ValidateCoupon(order.User_ID, orderreport.Product_ID, "", coupon_id)
			if valid == true {
				result := ur.DB.Table("orders").Where("order_id", order_id, "order_status", "pending").Update("coupon_id", coupon_id)
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
			valid, err := ur.ValidateCoupon(cart.User_ID, cart.Product_ID, cart_id, coupon_id)
			if valid == true {
				final_cart := domain.Final_Cart{}
				err := ur.DB.Where("cart_id", cart_id).First(&final_cart)
				if err.Error == nil {
					if final_cart.Coupon_ID != "not applied" {
						return errors.New("One Coupon already applied")
					}
					ur.Final_Cart(cart_id, coupon_id)
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
func (ur *UserRepo) CancelCoupon(cart_id, order_id, coupon_id string) error {
	if cart_id == "" {
		order := &domain.Order{}
		ur.DB.Where("order_id", order_id).First(&order)
		if order.Order_ID != "" {
			if order.Coupon_ID == "not applied" {
				return errors.New("not found applied coupons !")
			}

			result := ur.DB.Model(&domain.Order{}).Where("order_id", order_id).Update("coupon_id", "not applied")
			if result.Error != nil {
				return result.Error
			}

		} else {
			return errors.New("No product in the checklist")
		}
	} else {
		cart := &domain.Final_Cart{}
		ur.DB.Where("cart_id", cart_id).First(&cart)
		if cart.Cart_ID != "" {
			if cart.Coupon_ID == "not applied" {
				return errors.New("not found applied coupons !")
			}
			ur.Final_Cart(cart_id, "not applied")

		} else {
			return errors.New("No product in the cart")
		}
	}
	return nil
}

func (ur *UserRepo) Final_Cart(cart_id, coupon_id string) {
	cart := []domain.Cart{}
	final_cart := domain.Final_Cart{}
	result := ur.DB.Where("cart_id", cart_id).Find(&cart)
	if result.Error != nil {
		return
	}
	var totalPrice float32
	for i := range cart {
		totalPrice = totalPrice + cart[i].Total_Price
	}
	results := ur.DB.Where("cart_id", cart_id).First(&final_cart)
	if results.Error != nil {
		ur.DB.Create(&domain.Final_Cart{Cart_ID: cart_id, Total_Price: totalPrice, Coupon_ID: coupon_id, Discount: 0.0, Grand_Total: totalPrice})
		return
	}
	coupon := domain.Coupon{}
	if coupon_id != "not applied" {
		ur.DB.Where("coupon_code", coupon_id).First(&coupon)
		grand_total := final_cart.Grand_Total - coupon.Discount_amount
		ur.DB.Where("cart_id", cart_id).Updates(&domain.Final_Cart{Coupon_ID: coupon_id, Discount: coupon.Discount_amount, Grand_Total: grand_total})
		return
	} else {
		if final_cart.Coupon_ID != "not applied" {
			grand_total := final_cart.Grand_Total + final_cart.Discount
			ur.DB.Model(&domain.Final_Cart{}).Where("cart_id", cart_id).Updates(map[string]interface{}{"coupon_id": coupon_id, "grand_total": grand_total, "discount": 0.0})
			return
		}
		ur.DB.Where("cart_id", cart_id).Updates(&domain.Final_Cart{Total_Price: totalPrice, Grand_Total: totalPrice})
	}

}

//Shipping

func (ur *UserRepo) AddShippingDetails(user_id string, newAddress domain.ShippingDetails) error {

	address := &domain.ShippingDetails{}
	result := ur.DB.Where(&domain.ShippingDetails{Address: newAddress.Address, User_ID: user_id}).First(&address)
	if result.Error == nil {
		return errors.New("Entered input is already one of your shipping details")
	}
	result = ur.DB.Create(&domain.ShippingDetails{
		Shipping_ID: utils.GenerateID(),
		First_Name:  newAddress.First_Name,
		Last_Name:   newAddress.Last_Name,
		Email:       newAddress.Email,
		Phone:       newAddress.Phone,
		City:        newAddress.City,
		Street:      newAddress.Street,
		Address:     newAddress.Address,
		Pin_code:    newAddress.Pin_code,
		Land_Mark:   newAddress.Land_Mark,
		User_ID:     user_id})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (ur *UserRepo) ListShippingDetails(user_id string) ([]domain.ShippingDetailsResponse, error) {
	ShippingAddRes := []domain.ShippingDetailsResponse{}
	var totalRecords int64
	result := ur.DB.Table("shipping_details").Where(&domain.ShippingDetails{User_ID: user_id}).
		Select("shipping_id", "user_id", "concat(first_name,' ',last_name)as name", "email", "phone", "city", "street", "address", "pin_code", "land_mark").
		Where("deleted_at IS NULL").Find(&ShippingAddRes).Count(&totalRecords)

	if result.Error != nil {
		return nil, result.Error
	}
	if totalRecords == 0 {
		return nil, errors.New("No records found")
	}
	return ShippingAddRes, nil
}
func (ur *UserRepo) DeleteShippingDetails(user_id, address_id string) error {
	shipping_dtl := &domain.ShippingDetails{}
	fmt.Println(user_id)
	result := ur.DB.Where("shipping_id", address_id).Delete(&shipping_dtl)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (ur *UserRepo) CheckOut(cart_id, user_id, product_id, address_id string) (string, error) {
	result := ur.DB.Where("shipping_id", address_id).First(&domain.ShippingDetails{})
	if result.Error != nil {
		return "", result.Error
	}
	if cart_id != "" {
		final_cart := domain.Final_Cart{}
		cart := []domain.Cart{}
		result := ur.DB.Model(cart).Where("cart_id", cart_id).Find(&cart)
		if result.Error != nil {
			return "", result.Error
		}
		ur.DB.Where("cart_id", cart_id).First(&final_cart)
		id := utils.GenerateID()

		for i := range cart {
			result := ur.DB.Create(&domain.OrderReport{
				Order_ID:     id,
				Product_ID:   cart[i].Product_ID,
				Quantity:     cart[i].Quantity,
				TotalPrice:   cart[i].Total_Price,
				Order_Status: "pending",
			})
			if result.Error != nil {
				return "", result.Error
			}
		}
		result = ur.DB.Create(&domain.Order{
			Order_ID:       id,
			User_ID:        user_id,
			Shipping_ID:    address_id,
			Coupon_ID:      final_cart.Coupon_ID,
			Discount:       final_cart.Discount,
			Grand_Total:    final_cart.Grand_Total,
			GST:            (final_cart.Grand_Total * 12) / 100,
			Final:          (final_cart.Grand_Total) + ((final_cart.Grand_Total * 12) / 100),
			Order_Status:   "pending",
			Payment_Status: "pending",
		})
		if result.Error != nil {
			return "", result.Error
		}
		return id, nil
	} else {
		product, err := ur.ViewProduct(product_id)

		if err != nil {
			return "", err
		}
		id := utils.GenerateID()
		result := ur.DB.Create(&domain.OrderReport{
			Order_ID:     id,
			Product_ID:   product_id,
			Quantity:     1,
			TotalPrice:   *product.Unit_Price,
			Order_Status: "pending",
		})
		if result.Error != nil {
			return "", result.Error
		}
		result = ur.DB.Create(&domain.Order{
			Order_ID:       id,
			User_ID:        user_id,
			Shipping_ID:    address_id,
			Coupon_ID:      "not applied",
			Discount:       0,
			Grand_Total:    *product.Unit_Price,
			GST:            (*product.Unit_Price * 12) / 100,
			Final:          *product.Unit_Price + ((*product.Unit_Price * 12) / 100),
			Order_Status:   "pending",
			Payment_Status: "pending",
		})
		if result.Error != nil {
			return "", result.Error
		}
		return id, nil
	}
}
func (ur *UserRepo) OrderSummery(user_id string) (interface{}, domain.OrderSummery, error) {
	orderSummery := domain.OrderSummery{}
	type ProductDet struct {
		Product_code  string
		Product_Name  string
		Discription   string
		Product_Image *string
		Quantity      int
		TotalPrice    float32
	}
	orderreport := []domain.OrderReport{}
	productDet := []ProductDet{}
	orders := &domain.Order{}
	results := ur.DB.Table("orders").Where(&domain.Order{User_ID: user_id, Order_Status: "pending"}).Find(&orders)
	order_id := orders.Order_ID
	results = ur.DB.Where("order_id", order_id).Find(&orderreport)
	if results.Error != nil {
		return nil, orderSummery, results.Error
	}
	for i := range orderreport {
		product := &domain.Products{}
		ur.DB.Where(&domain.Products{Product_Code: orderreport[i].Product_ID}).First(&product)
		temp := ProductDet{
			Product_code:  product.Product_Code,
			Product_Name:  product.Product_Name,
			Discription:   product.Discription,
			Product_Image: product.Product_Image,
			Quantity:      orderreport[i].Quantity,
			TotalPrice:    orderreport[i].TotalPrice,
		}
		productDet = append(productDet, temp)
	}
	results = ur.DB.Model(&domain.Order{}).Where(&domain.Order{Order_ID: order_id, Order_Status: "pending"}).Joins("right join shipping_details on shipping_details.shipping_id=orders.shipping_id").
		Select("order_id", "orders.user_id", "shipping_details.shipping_id", "concat(first_name,' ',last_name)as shipping_name", "address as shipping_address",
			"coupon_id", "discount", "grand_total", "gst", "final", "mode_of_payment", "order_status", "payment_status").
		First(&orderSummery)
	if results.Error != nil {
		return nil, orderSummery, results.Error
	}
	return productDet, orderSummery, results.Error
}
func (ur *UserRepo) CancelOrder(order_id string) error {
	order := domain.Order{}
	result := ur.DB.Where(domain.Order{Order_ID: order_id, Order_Status: "pending"}).First(&order)
	if result.Error != nil {
		return result.Error
	}
	tx := ur.DB.Begin()

	if err := tx.Table("orders").Where("order_id", order_id).Updates(&domain.Order{Order_Status: "cancelled", Payment_Status: "not payed"}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Table("order_reports").Where("order_id", order_id).Update("order_status", "cancelled").Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
