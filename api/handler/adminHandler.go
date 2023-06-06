package handler

import (
	"ajalck/e_commerce/domain"
	services "ajalck/e_commerce/usecase/interface"
	"ajalck/e_commerce/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUseCase services.AdminUseCase
}

func NewAdminHandler(adminUseCase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: adminUseCase,
	}
}

// Registration

type NewAdmin struct {
	First_Name       string `json:"first_name" gorm:"not null" binding:"required,min=3"`
	Last_Name        string `json:"last_name"`
	Photo            string `json:"photo"`
	Email            string `json:"email" gorm:"not null" binding:"required,email"`
	Gender           string `json:"gender"`
	Phone            string `json:"phone" gorm:"not null" binding:"required,numeric,len=10"`
	Password         string `json:"password" gorm:"not null" binding:"required,min=6"`
	Conform_Password string `json:"conform_password" gorm:"not null" binding:"required,min=6"`
}

// @Summary Admin Signup
// @ID admin signup
// @Tags 1.Admin Registration
// @Param newAdmin body NewAdmin{} true "Register Admin"
// @Produce json
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/registration/signup [post]
func (ah *AdminHandler) CreateAdmin(c *gin.Context) {
	data := &NewAdmin{}
	if err := c.Bind(&data); err != nil || data.Password != data.Conform_Password {
		response := utils.ErrorResponse("Invalid inputs or missmatch in password !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusExpectationFailed)
		utils.ResponseJSON(c, response)
		return
	}
	newAdmin := domain.User{
		First_Name: data.First_Name,
		Last_Name:  data.Last_Name,
		Photo:      data.Photo,
		Email:      data.Email,
		Gender:     data.Gender,
		Phone:      data.Phone,
		Password:   data.Password,
	}
	err := ah.adminUseCase.CreateAdmin(newAdmin)
	if err != nil {
		response := utils.ErrorResponse("Couldnlt register a new admin !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	} else {
		response := utils.SuccessResponse("New admin registered successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// User Management

// @Summary List Users
// @ID list users
// @Tags 3.User Management
// @Security BearerAuth
// @Produce json
// @Param page query string true "Page No"
// @Param records query string true "No of records"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/userManagement/listusers/:page/:records [get]
func (ah *AdminHandler) ListUsers(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))

	users, metaData, err := ah.adminUseCase.ListUsers(page, perPage)
	type Page struct {
		Users    []domain.UserResponse
		MetaData utils.MetaData
	}
	result := Page{
		Users:    users,
		MetaData: metaData,
	}
	if err != nil {
		response := utils.ErrorResponse("Couldnlt list users !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is the users ...", result)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)

}

// @Summary View User
// @ID view user
// @Tags 3.User Management
// @Security BearerAuth
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/userManagement/viewuser/:id [get]
func (ah *AdminHandler) ViewUser(c *gin.Context) {

	id := c.Query("id")
	user, err := ah.adminUseCase.ViewUser(id)
	if err != nil {
		response := utils.ErrorResponse("Can't show user!", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is the user", user)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)

}

// @Summary Block User
// @ID block user
// @Tags 3.User Management
// @Security BearerAuth
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/userManagement/blockuser/:id [patch]
func (ah *AdminHandler) BlockUser(c *gin.Context) {

	id := c.Query("id")

	err := ah.adminUseCase.BlockUser(id)
	if err != nil {
		response := utils.ErrorResponse("Couldnt block user !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	} else {
		response := utils.SuccessResponse("Blocked user successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// @Summary Unblock User
// @ID unblock user
// @Tags 3.User Management
// @Security BearerAuth
// @Produce json
// @Param id query string true "User ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/userManagement/unblockuser/:id [patch]
func (ah *AdminHandler) UnblockUser(c *gin.Context) {

	id := c.Query("id")
	err := ah.adminUseCase.UnblockUser(id)
	if err != nil {
		response := utils.ErrorResponse("Couldnt unblock user !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	} else {
		response := utils.SuccessResponse("Unblocked user successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// @Summary List Blocked users
// @ID list blocked users
// @Tags 3.User Management
// @Security BearerAuth
// @Produce json
// @Param page query string true "Page No"
// @Param records query string true "No of records"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/userManagement/list/blockedusers/:page/:records [get]
func (ah *AdminHandler) ListBlockedUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	users, metaData, err := ah.adminUseCase.ListBlockedUsers(page, perPage)
	type Page struct {
		Users    []domain.UserResponse
		MetaData utils.MetaData
	}
	result := Page{
		Users:    users,
		MetaData: metaData,
	}
	if err != nil {
		response := utils.ErrorResponse("Couldnlt list users !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is the users ...", result)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)

}

// @Summary List Active users
// @ID list active users
// @Tags 3.User Management
// @Security BearerAuth
// @Produce json
// @Param page query string true "Page No"
// @Param records query string true "No of records"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/userManagement/list/activeusers/:page/:records [get]
func (ah *AdminHandler) ListActiveUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	users, metaData, err := ah.adminUseCase.ListActiveUsers(page, perPage)
	type Page struct {
		Users    []domain.UserResponse
		MetaData utils.MetaData
	}
	result := Page{
		Users:    users,
		MetaData: metaData,
	}
	if err != nil {
		response := utils.ErrorResponse("Couldnlt list users !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is the users ...", result)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)

}

// Category Management

// @Summary Add Categories
// @ID add categories
// @Tags 4.Category Management
// @Security BearerAuth
// @Produce json
// @Param newCateogory body domain.Category{} true "Add Category"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/categoryManagement/add [post]
func (ah *AdminHandler) AddCategory(c *gin.Context) {

	var NewCategory domain.Category
	if err := c.Bind(&NewCategory); err != nil {
		response := utils.ErrorResponse("Invalid inputs", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	err := ah.adminUseCase.AddCategory(NewCategory)
	if err != nil {
		response := utils.ErrorResponse("Couldnlt add category !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
	} else {
		response := utils.SuccessResponse("New Category added successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// @Summary list Categories
// @ID list categories
// @Tags 4.Category Management
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/categoryManagement/list [get]
func (ah *AdminHandler) ListCategory(c *gin.Context) {

	categories, err := ah.adminUseCase.ListCategory()
	if err != nil {
		response := utils.ErrorResponse("Couldnlt list category !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
	} else {
		response := utils.SuccessResponse("Here is the Categories", categories)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// @Summary Update Categories
// @ID update categories
// @Tags 4.Category Management
// @Security BearerAuth
// @Produce json
// @Param newCateogory body domain.Category{} true "update Category"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/categoryManagement/edit [patch]
func (ah *AdminHandler) EditCategory(c *gin.Context) {

	var NewCategory domain.Category
	if err := c.Bind(&NewCategory); err != nil {
		response := utils.ErrorResponse("Invalid inputs", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	err := ah.adminUseCase.EditCategory(NewCategory)

	if err != nil {
		response := utils.ErrorResponse("Couldnlt update category !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
	} else {
		response := utils.SuccessResponse("Category updated successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// @Summary Delete Categories
// @ID delete categories
// @Tags 4.Category Management
// @Security BearerAuth
// @Produce json
// @Param Cateogory body domain.Category{} true "Delete Category"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/categoryManagement/delete [delete]
func (ah *AdminHandler) DeleteCategory(c *gin.Context) {

	var category domain.Category
	if err := c.Bind(&category); err != nil {
		response := utils.ErrorResponse("Invalid inputs", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	err := ah.adminUseCase.DeleteCategory(category)
	if err != nil {
		response := utils.ErrorResponse("Couldnlt delete category !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
	} else {
		response := utils.SuccessResponse("Category deleted successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// Brand Management

// @Summary Add Brands
// @ID add brands
// @Tags 5.Brand Management
// @Security BearerAuth
// @Produce json
// @Param newBrand body domain.Brand{} true "Add Brand"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/brandManagement/add [post]
func (ah *AdminHandler) AddBrand(c *gin.Context) {

	var NewBrand domain.Brand
	if err := c.Bind(&NewBrand); err != nil {
		response := utils.ErrorResponse("Invalid inputs", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	err := ah.adminUseCase.AddBrand(NewBrand)
	if err != nil {
		response := utils.ErrorResponse("Couldn't add a new brand !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
	} else {
		response := utils.SuccessResponse("New brand added successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// @Summary List Brands
// @ID list brands
// @Tags 5.Brand Management
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/brandManagement/list [get]
func (ah *AdminHandler) ListBrands(c *gin.Context) {

	brands, err := ah.adminUseCase.ListBrands()
	if err != nil {
		response := utils.ErrorResponse("Couldn't list brands !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
	} else {
		response := utils.SuccessResponse("Here is the brands", brands)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// @Summary Update Brand
// @ID update brands
// @Tags 5.Brand Management
// @Security BearerAuth
// @Produce json
// @Param newBrand body domain.Brand{} true "Update Brand"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/brandManagement/edit [patch]
func (ah *AdminHandler) EditBrand(c *gin.Context) {

	var NewBrand domain.Brand
	if err := c.Bind(&NewBrand); err != nil {
		response := utils.ErrorResponse("Invalid inputs", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	err := ah.adminUseCase.EditBrand(NewBrand)
	if err != nil {
		response := utils.ErrorResponse("Couldn't update the brand !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
	} else {
		response := utils.SuccessResponse("brand updated successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// @Summary Delete Brand
// @ID delete brands
// @Tags 5.Brand Management
// @Security BearerAuth
// @Produce json
// @Param Brand body domain.Brand{} true "Delete Brand"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/brandManagement/delete [delete]
func (ah *AdminHandler) DeleteBrand(c *gin.Context) {

	var brand domain.Brand
	if err := c.Bind(&brand); err != nil {
		response := utils.ErrorResponse("Invalid inputs", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	err := ah.adminUseCase.DeleteBrand(brand)
	if err != nil {
		response := utils.ErrorResponse("Couldn't delete the brand !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
	} else {
		response := utils.SuccessResponse("brand deleted successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}
}

// Product Management

type NewProduct struct {
	Item          string   `json:"item" gorm:"not null"`
	Product_Name  string   `json:"product_name" gorm:"not null;unique"`
	Discription   string   `json:"discription" gorm:"not null"`
	Product_Image *string  `json:"product_image"`
	Category_id   uint     `json:"category_id" gorm:"not null"`
	Brand_id      uint     `json:"brand_id" gorm:"not null"`
	Size          *string  `json:"size" gorm:"not null"`
	Color         *string  `json:"color" gorm:"not null"`
	Unit_Price    *float32 `json:"unit_price" gorm:"not null"`
	Stock         *uint    `json:"stock" gorm:"not null"`
}

// @Summary Add Products
// @ID add products
// @Tags 6.Product Management
// @Security BearerAuth
// @Produce json
// @Param newProduct body NewProduct{} true "Add Product"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/productManagement/add [post]
func (ah *AdminHandler) AddProducts(c *gin.Context) {

	var newProduct *NewProduct
	if err := c.Bind(&newProduct); err != nil {
		response := utils.ErrorResponse("Invalid inputs !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	NewProducts := domain.Products{
		Item:          newProduct.Item,
		Product_Name:  newProduct.Product_Name,
		Discription:   newProduct.Discription,
		Product_Image: newProduct.Product_Image,
		Category_id:   newProduct.Category_id,
		Brand_id:      newProduct.Brand_id,
		Size:          newProduct.Size,
		Color:         newProduct.Color,
		Unit_Price:    newProduct.Unit_Price,
		Stock:         newProduct.Stock,
	}
	product_code, err := ah.adminUseCase.AddProducts(NewProducts)
	if err != nil {
		response := utils.ErrorResponse("Unable to add a new product !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Product added successfully", nil)
	c.Writer.Header().Set("product_code", product_code)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)

}

// @Summary List Products
// @ID list products
// @Tags 6.Product Management
// @Security BearerAuth
// @Produce json
// @Param page query string true "Page No"
// @Param records query string true "No of records"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/productManagement/list/:page/:records [get]
func (ah *AdminHandler) ListProducts(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))

	products, metaData, err := ah.adminUseCase.ListProducts(page, perPage)
	Page := struct {
		Products []domain.ProductResponse
		MetaData utils.MetaData
	}{Products: products,
		MetaData: metaData}
	if err != nil {
		response := utils.ErrorResponse("Couldnlt list products !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is the products ...", Page)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)

}

type UpdateProduct struct {
	Product_Code  string   `json:"product_code" gorm:"not null"`
	Product_Image *string  `json:"product_image"`
	Size          *string  `json:"size"`
	Color         *string  `json:"color"`
	Unit_Price    *float32 `json:"unit_price"`
	Stock         *uint    `json:"stock" `
}

// @Summary Update Products
// @ID update products
// @Tags 6.Product Management
// @Security BearerAuth
// @Produce json
// @Param newProduct body UpdateProduct{} true "Update Product"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/productManagement/edit [patch]
func (ah *AdminHandler) EditProducts(c *gin.Context) {

	var product *UpdateProduct
	if err := c.Bind(&product); err != nil {
		response := utils.ErrorResponse("Invalid inputs !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	NewProducts := domain.Products{
		Product_Code:  product.Product_Code,
		Product_Image: product.Product_Image,
		Size:          product.Size,
		Color:         product.Color,
		Unit_Price:    product.Unit_Price,
		Stock:         product.Stock,
	}
	err := ah.adminUseCase.EditProducts(NewProducts)
	if err != nil {
		response := utils.ErrorResponse("Unable update product !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	} else {
		response := utils.SuccessResponse("Product updated successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}

}

type ProductId struct {
	Product_Code string `json:"product_code"`
}

// @Summary Delete Products
// @ID delete products
// @Tags 6.Product Management
// @Security BearerAuth
// @Produce json
// @Param Productid body ProductId{} true "Product_ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/productManagement/delete [delete]
func (ah *AdminHandler) DeleteProducts(c *gin.Context) {
	product := &ProductId{}
	if err := c.Bind(&product); err != nil {
		response := utils.ErrorResponse("Invalid inputs !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	err := ah.adminUseCase.DeleteProducts(product.Product_Code)
	if err != nil {
		response := utils.ErrorResponse("Unable delete product !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)
		utils.ResponseJSON(c, response)
		return
	} else {
		response := utils.SuccessResponse("Product deleted successfully", nil)
		c.Writer.WriteHeader(http.StatusOK)
		utils.ResponseJSON(c, response)
	}

}

// Coupon
type AddCoupon struct {
	Coupon_Code     string  `json:"coupon_code"`
	Discount_amount float32 `json:"discount_amount"`
	User_ID         *string `json:"user_id"`
	Product_ID      *string `json:"product_id"`
	Min_Cost        float32 `json:"min_cost"`
	Validity        int64   `json:"validity"`
}

// @Summary Add Coupons
// @ID add coupons
// @Tags 7.Coupon Management
// @Security BearerAuth
// @Produce json
// @Param Coupon body AddCoupon{} true "Add Coupon"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/coupon/add [post]
func (ah *AdminHandler) AddCoupon(c *gin.Context) {
	addCoupon := AddCoupon{}
	if err := c.Bind(&addCoupon); err != nil {
		response := utils.ErrorResponse("Invalid inputs !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
	}
	newCoupon := domain.Coupon{
		Coupon_Code:     addCoupon.Coupon_Code,
		Discount_amount: addCoupon.Discount_amount,
		User_ID:         addCoupon.User_ID,
		Product_ID:      addCoupon.Product_ID,
		Min_Cost:        addCoupon.Min_Cost,
		Expires_At:      time.Now().AddDate(0, 0, int(addCoupon.Validity)),
	}
	fmt.Println(newCoupon.Expires_At)
	if newCoupon.Coupon_Code == "" {
		newCoupon.Coupon_Code = utils.GenerateID()
	}
	if err := ah.adminUseCase.AddCoupon(newCoupon); err != nil {
		response := utils.ErrorResponse("Couldn't add coupon", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusNotFound)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("New Coupon added", nil)
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(c, response)

}

// @Summary List Coupons on admin panel
// @ID list coupons on admin panel
// @Tags 7.Coupon Management
// @Security BearerAuth
// @Produce json
// @Param page query string true "Page No"
// @Param records query string true "No of records"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/coupon/list/:page/:records [get]
func (ah *AdminHandler) ListCoupon(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	coupon, metaData, err := ah.adminUseCase.ListCoupon(page, perPage)
	results := struct {
		Coupon   []domain.CouponResponse
		MetaData utils.MetaData
	}{
		Coupon:   coupon,
		MetaData: metaData,
	}
	if err != nil {
		response := utils.ErrorResponse("couldn't list coupons", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is the coupons", results)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

// @Summary Delete Coupons
// @ID delete coupons
// @Tags 7.Coupon Management
// @Security BearerAuth
// @Produce json
// @Param coupon_id query string true "Coupon_ID"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/coupon/delete/:couponid [delete]
func (ah *AdminHandler) DeleteCoupon(c *gin.Context) {
	coupon_id := c.Query("coupon_id")
	err := ah.adminUseCase.DeleteCoupon(coupon_id)
	if err != nil {
		response := utils.ErrorResponse("Couldn't remove the coupon", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("One coupon removed successfully", nil)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}

// @Summary Sales report
// @ID sales report
// @Tags Sales Report
// @Security BearerAuth
// @Produce json
// @Param page query string true "Page No"
// @Param records query string true "No of records"
// @Success 200 {object} utils.Response{}
// @Failure 422 {object} utils.Response{}
// @Router /admin/sales_report [get]
func (ah *AdminHandler) SalesReport(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	sales_report, metaData, err := ah.adminUseCase.SalesReport(page, perPage)
	if err != nil {
		response := utils.ErrorResponse("Couldn't list the sales report", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusNotFound)
		utils.ResponseJSON(c, response)
		return
	}
	Results := struct {
		Sales_Report interface{}
		MetaData     utils.MetaData
	}{
		Sales_Report: sales_report,
		MetaData:     metaData,
	}
	response := utils.SuccessResponse("Here is the sales report :", Results)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
}
