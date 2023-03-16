package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Email      string `json:"email"`
	Gender     string `json:"gender"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Status     string `json:"status"`
	User_Role  string `json:"user_role"`
	Level      string `json:"level"`
}
type WishList struct {
	Wishlist_ID uint `json:"category_id" gorm:"primarykey;unique;AUTO_INCREMENT"`
	User_ID     int  `json:"user_id"`
	Product_ID  int  `json:"product_id"`
}
type Cart struct {
	gorm.Model
	Cart_ID     string  `json:"cart_id" gorm:"not null"`
	User_ID     int     `json:"user_id"`
	Product_ID  int     `json:"product_id"`
	Coupon_id   int     `json:"coupon_id"`
	Quantity    int     `json:"quantity"`
	Unit_Price  float32 `json:"unit_price"`
	Total_Price float32 `json:"total_price"`
}
type ShippingDetails struct {
	gorm.Model
	First_Name string `json:"first_name" gorm:"not null"`
	Last_Name  string `json:"last_name"`
	Email      string `json:"email" gorm:"not null" binding:"required,email"`
	Phone      string `json:"phone" gorm:"not null" binding:"required,numeric,len=10"`
	City       string `json:"city" gorm:"not null"`
	Street     string `json:"street" gorm:"not null"`
	Address    string `json:"address" gorm:"not null"`
	Pin_code   string `json:"pin_code" gorm:"not null" binding:"required,numeric,len=6"`
	Land_Mark  string `json:"land_mark"`
	User_ID    uint   `json:"user_id"`
}
type Order struct {
	gorm.Model
	Order_ID        string  `json:"order_id" gorm:"not null"`
	User_ID         uint    `json:"user_id" gorm:"not null"`
	Product_ID      uint    `json:"product_id" gorm:"not null"`
	Shipping_ID     uint    `json:"shipping_id" gorm:"not null"`
	Coupon_ID       uint    `json:"coupon_id"`
	Payment_ID      uint    `json:"payment_id"`
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
