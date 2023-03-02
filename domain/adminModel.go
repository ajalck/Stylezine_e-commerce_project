package domain

import (
	"gorm.io/gorm"
)

type Products struct {
	gorm.Model
	Item          string   `json:"item"`
	Product_Name  string   `json:"product_name"`
	Discription   string   `json:"discription"`
	Product_Image string   `json:"product_image"`
	Category_id   uint     `json:"category_id" gorm:"not null"`
	Category      Category `json:"-" gorm:"foreignkey:Category_id;references:Category_ID"`
	Brand_id      uint     `json:"brand_id" gorm:"not null"`
	Brand         Brand    `json:"-" gorm:"foreignkey:Brand_id;references:Brand_ID"`
	Wishlist_id   uint     `json:"wishlist_id"`
	WishList      WishList `json:"-" gorm:"foreignkey:Wishlist_id;references:Wishlist_ID"`
	Cart_id       uint     `json:"cart_id"`
	Cart          Cart     `json:"-" gorm:"foreignkey:Cart_id;references:Cart_ID"`
	Size          string   `json:"size" gorm:"not null"`
	Color         string   `json:"color" gorm:"not null"`
	Unit_Price    float32  `json:"unit_price"`
	Stock         uint     `json:"stock"`
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
