package handler

import (
	"ajalck/e_commerce/domain"
	services "ajalck/e_commerce/usecase/interface"
	"ajalck/e_commerce/utils"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

func NewUserHandler(userUseCase services.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

type NewUser struct {
	First_Name       string `json:"first_name" gorm:"not null" binding:"required,min=3"`
	Last_Name        string `json:"last_name"`
	Photo            string `json:"photo"`
	Email            string `json:"email" gorm:"not null" binding:"required,email"`
	Gender           string `json:"gender"`
	Phone            string `json:"phone" gorm:"not null" binding:"required,numeric,len=10"`
	Password         string `json:"password" gorm:"not null" binding:"required,min=6"`
	Conform_Password string `json:"conform_password" gorm:"not null" binding:"required,min=6"`
}

// @title Go + Gin Stylezine API
// @version 1.0
// @description This is a sample server Job Portal server. You can visit the GitHub repository at https://github.com/ajalck/Stylezine_e-commerce_project

// @contact.name API Support
// @contact.url https://github.com/ajalck/ajal_portfolio
// @contact.email ack6627@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:5050
// @BasePath
// @query.collection.format multi

// @Summary Add user to database
// @ID create user
// @Tags 10.User Registration
// @Produce json
// @Param newUser body NewUser{} true "New User"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/signup [post]
func (uh *UserHandler) CreateUser(c *gin.Context) {

	data := NewUser{}
	if err := c.Bind(&data); err != nil || data.Password != data.Conform_Password {
		response := utils.ErrorResponse("Invalid inputs or missmatch in password !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusExpectationFailed)
		utils.ResponseJSON(c, response)
		return
	}
	newUser := domain.User{
		First_Name: data.First_Name,
		Last_Name:  data.Last_Name,
		Photo:      data.Photo,
		Email:      data.Email,
		Gender:     data.Gender,
		Phone:      data.Phone,
		Password:   data.Password,
	}
	err := uh.userUseCase.CreateUser(newUser)
	if err != nil {
		response := utils.ErrorResponse("Couldnlt register a new user !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		utils.ResponseJSON(c, response)
		return
	} else {
		response := utils.SuccessResponse("New user registered successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// @Summary List products to user
// @ID list products to user
// @Tags List Products
// @Produce json
// @Param page query string true "Page No"
// @Param records query string true "No of records"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/listproducts/:page/:records [get]
func (uh *UserHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	users, metaData, err := uh.userUseCase.ListProducts(page, perPage)
	type Page struct {
		users    []domain.ProductResponse
		metaData utils.MetaData
	}
	result := Page{
		users:    users,
		metaData: metaData,
	}
	if err == nil {
		c.JSON(200, result.users)
		c.JSON(http.StatusFound, result.metaData)
		return
	}
	if err != nil {
		c.JSON(400, err.Error())
	}
}

// @Summary Add product to wishlist
// @ID user add wishlist
// @Tags User Wishlist
// @Security BearerAuth
// @Produce json
// @Param product_id query string true "Product_ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/wishlist/add/:productid [post]
func (uh *UserHandler) AddWishlist(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	product_id, _ := strconv.Atoi(c.Query("product_id"))
	err := uh.userUseCase.AddWishlist(user_id, product_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't add new item to wishlist", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("New item added to wishlist", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

// @Summary View wishlist
// @ID user view wishlist
// @Tags User Wishlist
// @Security BearerAuth
// @Produce json
// @Param page query string true "page no"
// @Param records query string true "no of records"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/wishlist/view/:page/:records [get]
func (uh *UserHandler) ViewWishList(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	wishList, metaData, err := uh.userUseCase.ViewWishList(user_id, page, perPage)
	var results struct {
		WishLists []domain.WishListResponse
		MetaData  utils.MetaData
	}
	results.WishLists = wishList
	results.MetaData = metaData
	if err != nil {
		response := utils.ErrorResponse("couldn't reach to your wishlist", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is your wishlist...", results)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

// @Summary Remove product from wishlist
// @ID user delete wishlist
// @Tags User Wishlist
// @Security BearerAuth
// @Produce json
// @Param product_id query string true "Product_ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/wishlist/remove/:productid [delete]
func (uh *UserHandler) DeleteWishList(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	product_id, _ := strconv.Atoi(c.Query("product_id"))
	err := uh.userUseCase.DeleteWishList(user_id, product_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't remove the item from your wishlist", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("one item removed successfully from your wishlist", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

//Cart

// @Summary Add product to cart
// @ID user add cart
// @Tags User Cart
// @Security BearerAuth
// @Produce json
// @Param product_id query string true "Product_ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/cart/add/:productid [post]
func (uh *UserHandler) AddCart(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	product_id, _ := strconv.Atoi(c.Query("product_id"))
	err, cart_id := uh.userUseCase.AddCart(user_id, product_id)
	c.Writer.Header().Set("cart_id", cart_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't add new item to cart", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("New item added to cart", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

// @Summary View cart
// @ID user view cart
// @Tags User Cart
// @Security BearerAuth
// @Produce json
// @Param page query string true "page no"
// @Param records query string true "no of records"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/cart/view/:page/:records [get]
func (uh *UserHandler) ViewCart(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	cart, metaData, err := uh.userUseCase.ViewCart(user_id, page, perPage)
	results := struct {
		cart     []domain.CartResponse
		MetaData utils.MetaData
	}{
		cart:     cart,
		MetaData: metaData,
	}
	if err != nil {
		response := utils.ErrorResponse("couldn't reach to your cart", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is your cart...", results.cart)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
	c.JSON(200, results.MetaData)
}

// @Summary remove product from cart
// @ID user delete cart
// @Tags User Cart
// @Security BearerAuth
// @Produce json
// @Param product_id query string true "Product_ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/cart/remove/:productid [post]
func (uh *UserHandler) DeleteCart(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	product_id, _ := strconv.Atoi(c.Query("product_id"))
	err := uh.userUseCase.DeleteCart(user_id, product_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't remove an item from your cart", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("One item removed successfully from your cart", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

//Coupon

// @Summary List coupon
// @ID user list coupon
// @Tags User Coupon
// @Security BearerAuth
// @Produce json
// @Param product_id query string true "Product_ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/coupon/listcoupon/:productid [get]
func (uh *UserHandler) ListCoupon(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	product_id, _ := strconv.Atoi(c.Query("product_id"))

	coupons, err := uh.userUseCase.ListCoupon(user_id, product_id)
	if err != nil {
		response := utils.ErrorResponse("No coupons found !", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("here is the coupons", coupons)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

// @Summary Apply coupon
// @ID user apply coupon
// @Tags User Coupon
// @Security BearerAuth
// @Produce json
// @Param cart_id query string true "cart_ID"
// @Param order_id query string true "order_id"
// @Param coupon_id query string true "coupon_id"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/coupon/applycoupon/:cartid/:orderid/:couponid [post]
func (uh *UserHandler) ApplyCoupon(c *gin.Context) {
	cart_id := c.Query("cart_id")
	order_id := c.Query("order_id")
	coupon_id, _ := strconv.Atoi(c.Query("coupon_id"))
	err := uh.userUseCase.ApplyCoupon(cart_id, order_id, coupon_id)
	if err != nil {
		response := utils.ErrorResponse("Coupon couldn't applied !", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Coupon applied successfully", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

// @Summary Cancel coupon
// @ID user cancel coupon
// @Tags User Coupon
// @Security BearerAuth
// @Produce json
// @Param cart_id query string true "cart_ID"
// @Param order_id query string true "order_id"
// @Param coupon_id query string true "coupon_id"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/coupon/cancelcoupon/:cartid/:orderid/:couponid [delete]
func (uh *UserHandler) CancelCoupon(c *gin.Context) {
	cart_id := c.Query("cart_id")
	order_id := c.Query("order_id")
	coupon_id, _ := strconv.Atoi(c.Query("coupon_id"))
	err := uh.userUseCase.CancelCoupon(cart_id, order_id, coupon_id)
	if err != nil {
		response := utils.ErrorResponse("Coupon couldn't cancelled !", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Coupon cancelled successfully", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

//Shipping

// @Summary Add Shipping details
// @ID user add shipping details
// @Tags User Shipping details
// @Security BearerAuth
// @Produce json
// @Param newAddress body domain.ShippingDetails{} true "Shipping details"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/shipping/adddetails [post]
func (uh *UserHandler) AddShippingDetails(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	var newAddress domain.ShippingDetails
	if err := c.Bind(&newAddress); err != nil {
		response := utils.ErrorResponse("Invalid inputs", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	err := uh.userUseCase.AddShippingDetails(user_id, newAddress)
	if err != nil {
		response := utils.ErrorResponse("Failed to add new shipping details", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusConflict)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("New shipping details added successfully", nil)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)
}

// @Summary List Shipping details
// @ID user list shipping details
// @Tags User Shipping details
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/shipping/listdetails [get]
func (uh *UserHandler) ListShippingDetails(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	shippingDetails, err := uh.userUseCase.ListShippingDetails(user_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't list shipping details", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusNotFound)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is your shipping details", shippingDetails)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)
}

// @Summary Delete Shipping details
// @ID user delete shipping details
// @Tags User Shipping details
// @Security BearerAuth
// @Produce json
// @Param address_id query string true "address_ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/shipping/deletedetails/:addressid [delete]
func (uh *UserHandler) DeleteShippingDetails(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	address_id, _ := strconv.Atoi(c.Query("address_id"))
	err := uh.userUseCase.DeleteShippingDetails(user_id, address_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't delete shipping details", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusNotFound)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("successfully removed selected shipping details", nil)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)
}

//Order

// @Summary Add to checkout
// @ID user add to checkout
// @Tags User Order
// @Security BearerAuth
// @Produce json
// @Param cart_id query string true "cart_ID"
// @Param product_id query string true "product_ID"
// @Param address_id query string true "address_ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/order/checkout/:cartid/:productid/:shippingid [post]
func (uh *UserHandler) CheckOut(c *gin.Context) {
	user_id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	cart_id := c.Query("cart_id")
	product_id, _ := strconv.Atoi(c.Query("product_id"))
	address_id, _ := strconv.Atoi(c.Query("address_id"))

	id, err := uh.userUseCase.CheckOut(cart_id, user_id, product_id, address_id)
	if err != nil {
		response := utils.ErrorResponse("Failed to add to checkout !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusNotFound)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("added to checkout", nil)
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("order_id", id)
	utils.ResponseJSON(c, response)

}

// @Summary View order summery
// @ID user view order summery
// @Tags User Order
// @Security BearerAuth
// @Produce json
// @Param order_id query string true "order_ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/order/ordersummery/:orderid [get]
func (uh *UserHandler) OrderSummery(c *gin.Context) {
	order_id := c.Query("order_id")
	orderSummery, err := uh.userUseCase.OrderSummery(order_id)

	if err != nil {
		response := utils.ErrorResponse("Couldn't display order summery", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusNotFound)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Order summery :", orderSummery)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)

}

// @Summary update order
// @ID user update order
// @Tags User Order
// @Security BearerAuth
// @Produce json
// @Param order_id query string true "order_ID"
// @Param OrderUpdates body interface{} true "Order Updates"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /user/order/update/:orderid [patch]
func (uh *UserHandler) UpdateOrder(c *gin.Context) {
	type OrderUpdates struct {
		Quantity        int    `json:"quantity"`
		Mode_of_Payment string `json:"mode_of_payment"`
	}
	updates := &OrderUpdates{}
	if err := c.Bind(&updates); err != nil {
		response := utils.ErrorResponse("Couldn't update order", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusNotFound)
		utils.ResponseJSON(c, response)
	}
	order_id := c.Query("order_id")
	product_id := c.Query("product_id")
	err := uh.userUseCase.UpdateOrder(order_id, product_id, updates)
	if err != nil {

		return
	}
	response := utils.SuccessResponse("order details updated :", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}
