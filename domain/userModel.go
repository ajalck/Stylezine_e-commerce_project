package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	Status       string `json:"status"`
	User_Role    string `json:"user_role"`
	Verification bool   `json:"verification"`
	Token        string `json:"token"`
}
type WishList struct {
	Wishlist_ID uint `json:"category_id" gorm:"primarykey;unique;AUTO_INCREMENT"`
	User_ID     int  `json:"user_id"`
	Product_ID  int  `json:"product_id"`
}
