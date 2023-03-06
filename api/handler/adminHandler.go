package handler

import (
	"ajalck/e_commerce/domain"
	services "ajalck/e_commerce/usecase/interface"
	"ajalck/e_commerce/utils"
	"net/http"
	"strconv"

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

func (ah *AdminHandler) CreateAdmin(c *gin.Context) {
	var newAdmin domain.User
	if err := c.Bind(&newAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Signup Inputs"})
		return
	}
	err := ah.adminUseCase.CreateAdmin(c, newAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "couldn't create a new admin"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "New admin created successfully"})
	}
}

// User Management

func (ah *AdminHandler) ListUsers(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))

	users, metaData, err := ah.adminUseCase.ListUsers(page, perPage)
	type Page struct {
		users    []domain.UserResponse
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
func (ah *AdminHandler) ViewUser(c *gin.Context) {

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, err)
		c.Abort()
	}
	user, err := ah.adminUseCase.ViewUser(id)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	c.JSON(200, user)

}
func (ah *AdminHandler) BlockUser(c *gin.Context) {

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, "Invalid input")
		return
	}
	err = ah.adminUseCase.BlockUser(id)
	if err != nil {
		c.JSON(400, err.Error())
	} else {
		c.JSON(200, "User is blocked")
	}
}
func (ah *AdminHandler) UnblockUser(c *gin.Context) {

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, "Invalid input")
		return
	}
	err = ah.adminUseCase.UnblockUser(id)
	if err != nil {
		c.JSON(400, err.Error())
	} else {
		c.JSON(200, "User is unblocked")
	}
}
func (ah *AdminHandler) ListBlockedUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	users, metaData, err := ah.adminUseCase.ListBlockedUsers(page, perPage)
	type Page struct {
		users    []domain.UserResponse
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
func (ah *AdminHandler) ListActiveUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	users, metaData, err := ah.adminUseCase.ListActiveUsers(page, perPage)
	type Page struct {
		users    []domain.UserResponse
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

// Category Management

func (ah *AdminHandler) AddCategory(c *gin.Context) {

	var NewCategory domain.Category
	if err := c.Bind(&NewCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inputs !"})
		return
	}
	err := ah.adminUseCase.AddCategory(NewCategory)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to add a new category !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "New Category added successfully"})
	}
}
func (ah *AdminHandler) EditCategory(c *gin.Context) {

	var NewCategory domain.Category
	if err := c.Bind(&NewCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inputs !"})
		return
	}
	err := ah.adminUseCase.EditCategory(NewCategory)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to update category !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "Category updated successfully"})
	}
}
func (ah *AdminHandler) DeleteCategory(c *gin.Context) {

	var category domain.Category
	if err := c.Bind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inputs !"})
		return
	}
	err := ah.adminUseCase.DeleteCategory(category)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to delete category !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "Category deleted successfully"})
	}
}

// Brand Management

func (ah *AdminHandler) AddBrand(c *gin.Context) {

	var NewBrand domain.Brand
	if err := c.Bind(&NewBrand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Inputs !"})
		return
	}
	err := ah.adminUseCase.AddBrand(NewBrand)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to add a new brand !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "New Brand added successfully"})
	}
}
func (ah *AdminHandler) EditBrand(c *gin.Context) {

	var NewBrand domain.Brand
	if err := c.Bind(&NewBrand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Inputs !"})
		return
	}
	err := ah.adminUseCase.EditBrand(NewBrand)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to update brand !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "Brand updated successfully"})
	}
}
func (ah *AdminHandler) DeleteBrand(c *gin.Context) {

	var brand domain.Brand
	if err := c.Bind(&brand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Inputs !"})
		return
	}
	err := ah.adminUseCase.DeleteBrand(brand)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to delete brand !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "Brand deleted successfully"})
	}
}

// Product Management

func (ah *AdminHandler) AddProducts(c *gin.Context) {

	var NewProducts domain.Products
	if err := c.Bind(&NewProducts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inputs !"})
		return
	}
	err := ah.adminUseCase.AddProducts(NewProducts)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to add a new product !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "New Product added successfully "})
	}

}
func (ah *AdminHandler) EditProducts(c *gin.Context) {

	var NewProducts domain.Products
	if err := c.Bind(&NewProducts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inputs !"})
		return
	}
	err := ah.adminUseCase.EditProducts(NewProducts)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to update product !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "Product updated successfully "})
	}

}
func (ah *AdminHandler) DeleteProducts(c *gin.Context) {

	var product domain.Products
	if err := c.Bind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Inputs !"})
		return
	}
	err := ah.adminUseCase.DeleteProducts(product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to delete product !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "Product deleted successfully "})
	}

}

// Coupon
func (ah *AdminHandler) AddCoupon(c *gin.Context) {
	newCoupon := domain.Coupon{}
	if err := c.Bind(&newCoupon); err != nil {
		response := utils.ErrorResponse("Invalid inputs !", err.Error(), nil)
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(c, response)
		return
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
	return

}
func (ah *AdminHandler) ListCoupon(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("records"))
	coupon, metaData, err := ah.adminUseCase.ListCoupon(page, perPage)
	results := struct {
		coupon   []domain.CouponResponse
		MetaData utils.MetaData
	}{
		coupon:   coupon,
		MetaData: metaData,
	}
	if err != nil {
		response := utils.ErrorResponse("couldn't list coupons", err.Error(), nil)
		c.Writer.WriteHeader(400)
		utils.ResponseJSON(c, response)
		return
	}
	response := utils.SuccessResponse("Here is the coupons", results.coupon)
	c.Writer.WriteHeader(200)
	utils.ResponseJSON(c, response)
	c.JSON(200, results.MetaData)
}
func (ah *AdminHandler) DeleteCoupon(c *gin.Context) {
	coupon_id, _ := strconv.Atoi(c.Query("coupon_id"))
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
