package domain

import (
	"time"

	"gorm.io/gorm"
)

type Products struct {
	gorm.Model
	Product_Code  string   `json:"product_code" gorm:"not null"`
	Item          string   `json:"item" gorm:"not null"`
	Product_Name  string   `json:"product_name" gorm:"not null;unique"`
	Discription   string   `json:"discription" gorm:"not null"`
	Product_Image *string  `json:"product_image"`
	Category_id   uint     `json:"category_id" gorm:"not null"`
	Category      Category `json:"-" gorm:"foreignkey:Category_id;references:Category_ID"`
	Brand_id      uint     `json:"brand_id" gorm:"not null"`
	Brand         Brand    `json:"-" gorm:"foreignkey:Brand_id;references:Brand_ID"`
	Size          *string  `json:"size"`
	Color         *string  `json:"color"`
	Unit_Price    *float32 `json:"unit_price"`
	Stock         *uint    `json:"stock"`
	Rating        float32  `json:"rating"`
	Status        string   `json:"status"`
}
type Category struct {
	Category_ID   uint   `json:"category_id" gorm:"primarykey;unique;AUTO_INCREMENT"`
	Category_name string `json:"category_name" gorm:"unique;"`
}
type Brand struct {
	Brand_ID          uint   `json:"brand_id" gorm:"primarykey;autoIncrement:true;unique"`
	Brand_Name        string `json:"brand_name" gorm:"unique"`
	Brand_Discription string `json:"brand_discription"`
}
type Coupon struct {
	gorm.Model
	Coupon_Code     string    `json:"coupon_code"`
	Discount_amount float32   `json:"discount_amount"`
	User_ID         *string   `json:"user_id"`
	Product_ID      *string   `json:"product_id"`
	Min_Cost        float32   `json:"min_cost"`
	Expires_At      time.Time `json:"expires_at"`
	Coupon_Status   string    `json:"coupon_status"`
}
type OrderReport struct {
	Order_ID     string  `json:"order_id" gorm:"not null"`
	Order        Order   `gorm:"foreignkey:Order_ID;references:Order_ID"`
	Product_ID   string  `json:"product_id" gorm:"not null"`
	Quantity     int     `json:"quantity" gorm:"not null"`
	TotalPrice   float32 `json:"totalprice"`
	Order_Status string  `json:"order_status"`
}
