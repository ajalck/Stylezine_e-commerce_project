package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	User_ID    string `json:"user_id" gorm:"unique"`
	First_Name string `json:"first_name" gorm:"not null" binding:"required,min=3"`
	Last_Name  string `json:"last_name"`
	Photo      string `json:"photo"`
	Email      string `json:"email" gorm:"not null" binding:"required,email"`
	Gender     string `json:"gender"`
	Phone      string `json:"phone" gorm:"not null" binding:"required,numeric,len=10"`
	Password   string `json:"password" gorm:"not null" binding:"required,min=6"`
	Status     string `json:"status"`
	User_Role  string `json:"user_role"`
	Level      string `json:"level"`
}
type WishList struct {
	gorm.Model
	Wishlist_ID string `json:"wishlist_id" gorm:"primarykey;unique;AUTO_INCREMENT"`
	User_ID     string `json:"user_id"`
	Product_ID  string `json:"product_id"`
}
type Cart struct {
	gorm.Model
	Cart_ID     string  `json:"cart_id" gorm:"not null"`
	User_ID     string  `json:"user_id"`
	Product_ID  string  `json:"product_id"`
	Coupon_id   string  `json:"coupon_id"`
	Quantity    int     `json:"quantity"`
	Unit_Price  float32 `json:"unit_price"`
	Total_Price float32 `json:"total_price"`
}
type ShippingDetails struct {
	gorm.Model
	Shipping_ID string `json:"shipping_id"`
	First_Name  string `json:"first_name" gorm:"not null"`
	Last_Name   string `json:"last_name"`
	Email       string `json:"email" gorm:"not null" binding:"required,email"`
	Phone       string `json:"phone" gorm:"not null" binding:"required,numeric,len=10"`
	City        string `json:"city" gorm:"not null"`
	Street      string `json:"street" gorm:"not null"`
	Address     string `json:"address" gorm:"not null"`
	Pin_code    string `json:"pin_code" gorm:"not null" binding:"required,numeric,len=6"`
	Land_Mark   string `json:"land_mark"`
	User_ID     string `json:"user_id"`
}
type Order struct {
	gorm.Model
	Order_ID        string  `json:"order_id" gorm:"not null"`
	User_ID         string  `json:"user_id" gorm:"not null"`
	Product_ID      string  `json:"product_id" gorm:"not null"`
	Shipping_ID     string  `json:"shipping_id" gorm:"not null"`
	Coupon_ID       string  `json:"coupon_id"`
	Payment_ID      string  `json:"payment_id"`
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
