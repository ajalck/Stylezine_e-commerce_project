package config

import (
	"ajalck/e_commerce/domain"
	"log"

	"gorm.io/gorm"
)

func SyncDB(db *gorm.DB) {

	if err := (db.AutoMigrate(&domain.User{})); err != nil {
		log.Println("Failed to sync 'User' table")
	}
	if err := (db.AutoMigrate(&domain.SignedDetails{})); err != nil {
		log.Println("Failed to sync 'SignedDetails' table")
	}
	if err := (db.AutoMigrate(&domain.Products{})); err != nil {
		log.Println("Failed to sync 'Products' table")
	}
	if err := (db.AutoMigrate(&domain.Category{})); err != nil {
		log.Println("Failed to sync 'Category' table")
	}
	if err := (db.AutoMigrate(&domain.Brand{})); err != nil {
		log.Println("Failed to sync 'Brand' table")
	}
	if err := (db.AutoMigrate(&domain.WishList{})); err != nil {
		log.Println("Failed to sync 'WishList' table")
	}
	if err := (db.AutoMigrate(&domain.Cart{})); err != nil {
		log.Println("Failed to sync 'Cart' table")
	}
	if err := (db.AutoMigrate(&domain.ShippingDetails{})); err != nil {
		log.Println("Failed to sync 'ShippingDetails' table")
	}
	if err := (db.AutoMigrate(&domain.Coupon{})); err != nil {
		log.Println("Failed to sync 'Coupon' table")
	}
	if err := (db.AutoMigrate(&domain.Order{})); err != nil {
		log.Println("Failed to sync 'Order' table")
	}
	if err := (db.AutoMigrate(&domain.Final_Cart{})); err != nil {
		log.Println("Failed to sync 'Final_Cart' table")
	}
	if err := (db.AutoMigrate(&domain.Final_Cart{})); err != nil {
		log.Println("Failed to sync 'Final_Cart' table")
	}
	if err := (db.AutoMigrate(&domain.OrderReport{})); err != nil {
		log.Println("Failed to sync 'OrderReport' table")
	}
}
