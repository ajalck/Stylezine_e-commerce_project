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

}
