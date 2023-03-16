package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserResponse struct {
	ID           int    `json:"id"`
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
	Phone        string `json:"phone"`
	Status       string `json:"status"`
	User_Role    string `json:"user_role"`
	Verification bool   `json:"verification"`
}
type ProductResponse struct {
	ID            int     `json:"id"`
	Item          string  `json:"item"`
	Product_Name  string  `json:"product_name"`
	Discription   string  `json:"discription"`
	Product_Image string  `json:"product_image"`
	Category_name string  `json:"category_name"`
	Brand_name    string  `json:"brand_name"`
	Size          string  `json:"size" gorm:"not null"`
	Color         string  `json:"color" gorm:"not null"`
	Unit_Price    float32 `json:"unit_price"`
	Stock         uint    `json:"stock"`
	Rating        float32 `json:"rating"`
	Status        string  `json:"status"`
}
type SignedDetails struct {
	UserId   int
	Username string
	UserRole string
	jwt.StandardClaims
}
type WishListResponse struct {
	ID            int    `json:"id"`
	Item          string `json:"item"`
	Product_Name  string `json:"product_name"`
	Product_Image string `json:"product_image"`
	Size          string `json:"size" gorm:"not null"`
	Color         string `json:"color" gorm:"not null"`
	Status        string `json:"status"`
}
type CartResponse struct {
	User_id       int     `json:"user_id"`
	Product_id    int     `json:"product_id"`
	Item          string  `json:"item"`
	Product_Name  string  `json:"product_name"`
	Product_Image string  `json:"product_image"`
	Size          string  `json:"size" gorm:"not null"`
	Color         string  `json:"color" gorm:"not null"`
	Count         int     `json:"count"`
	TotalPrice    float32 `json:"total_price"`
	Status        string  `json:"status"`
}
type ShippingDetailsResponse struct {
	ID        int    `json:"id"`
	User_ID   int    `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	City      string `json:"city"`
	Street    string `json:"street"`
	Address   string `json:"address"`
	Pin_code  string `json:"pin_code"`
	Land_Mark string `json:"land_mark"`
}

type CouponResponse struct {
	ID              uint      `json:"coupon_id"`
	Coupon_Code     string    `json:"coupon_code"`
	Discount_amount float32   `json:"discount_amount"`
	User_ID         uint      `json:"user_id"`
	Product_ID      uint      `json:"product_id"`
	Min_Cost        float32   `json:"min_cost"`
	Expires_At      time.Time `json:"expires_at"`
	Coupon_Status   string    `json:"coupon_status"`
}
type OrderSummery struct {
	// User_ID          uint    `json:"user_id" gorm:"not null"`
	// Shipping_Name    string  `json:"shipping_name"`
	// Shipping_Address string  `json:"shipping_address"`
	// Product_Name     string  `json:"product_name"`
	// Discription      string  `json:"discription"`
	// Product_Image    string  `json:"product_image"`
	User_ID         uint    `json:"user_id" gorm:"not null"`
	Product_ID      uint    `json:"product_id" gorm:"not null"`
	Shipping_ID     uint    `json:"shipping_id" gorm:"not null"`
	Coupon_ID       uint    `json:"coupon_id"`
	Quantity        int     `json:"quantity" gorm:"not null"`
	Discount        float32 `json:"discount"`
	TotalPrice      float32 `json:"totalprice"`
	Grand_Total     float32 `json:"grand_total"`
	GST             float32 `json:"gst"`
	Final           float32 `json:"final"`
	Mode_of_Payment string  `json:"mode_of_payment"`
	Order_Status    string  `json:"order_status"`
	Payment_Status  string  `json:"payment_status"`
}
