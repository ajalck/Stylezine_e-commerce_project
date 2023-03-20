package usecase

import (
	"ajalck/e_commerce/domain"
	repoInt "ajalck/e_commerce/repository/interface"
	services "ajalck/e_commerce/usecase/interface"
	"ajalck/e_commerce/utils"
	"errors"
	"fmt"
)

type adminUseCase struct {
	adminRepo repoInt.AdminRepository
}

func NewAdminUseCase(repo repoInt.AdminRepository) services.AdminUseCase {
	return &adminUseCase{
		adminRepo: repo,
	}
}
func (au *adminUseCase) CreateAdmin(newAdmin domain.User) error {

	newAdmin.User_Role = "admin"

	if _, err := au.adminRepo.FindAdmin(newAdmin.Email, newAdmin.User_Role); err == nil {
		return errors.New("Admin already exists")
	}
	//Hashing password

	newAdmin.Password = HashPassword(newAdmin.Password)
	newAdmin.User_ID = utils.GenerateID()
	newAdmin.Status = "active"
	err := au.adminRepo.CreateAdmin(newAdmin)
	if err != nil {
		return err
	}
	return nil
}

// User Management

func (au *adminUseCase) ListUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error) {

	users, metaData, err := au.adminRepo.ListUsers(page, perPage)
	if err != nil {
		return users, metaData, err
	}
	return users, metaData, nil

}
func (au *adminUseCase) ViewUser(id string) (domain.UserResponse, error) {
	user, err := au.adminRepo.ViewUser(id)
	if err != nil {
		return user, err
	}
	return user, nil

}

func (au *adminUseCase) BlockUser(id string) error {
	user, err := au.adminRepo.ViewUser(id)
	if err != nil {
		return err
	}
	if user.Status == "blocked" {
		err := errors.New("user is already blocked")
		return err
	}
	au.adminRepo.BlockUser(id)
	return nil
}
func (au *adminUseCase) UnblockUser(id string) error {
	user, err := au.adminRepo.ViewUser(id)
	if err != nil {
		return err
	}
	if user.Status == "active" {
		err := errors.New("user is already active")
		fmt.Println(err)
		return err
	}
	au.adminRepo.UnblockUser(id)
	return nil
}
func (au *adminUseCase) ListBlockedUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error) {

	users, metaData, err := au.adminRepo.ListBlockedUsers(page, perPage)
	if err != nil {
		return users, metaData, err
	}
	return users, metaData, nil
}
func (au *adminUseCase) ListActiveUsers(page, perPage int) ([]domain.UserResponse, utils.MetaData, error) {

	users, metaData, err := au.adminRepo.ListActiveUsers(page, perPage)
	if err != nil {
		return users, metaData, err
	}
	return users, metaData, nil
}

// Category Management

func (au *adminUseCase) AddCategory(NewCategory domain.Category) error {
	err := au.adminRepo.AddCategory(NewCategory)
	if err != nil {
		return err
	}
	return nil

}
func (au *adminUseCase) ListCategory() ([]domain.Category, error) {

	categories, err := au.adminRepo.ListCategory()
	if err != nil {
		return nil, err
	}
	return categories, nil

}
func (au *adminUseCase) EditCategory(NewCategory domain.Category) error {

	category, err := au.adminRepo.ViewCategory(NewCategory)

	if err != nil {
		return err
	} else {
		if category.Category_name != NewCategory.Category_name {
			if err := au.adminRepo.EditCategory(NewCategory); err != nil {
				return err
			}
			return nil
		}
		return errors.New("entered category name is same as the existing")
	}
}
func (au *adminUseCase) DeleteCategory(category domain.Category) error {

	category, err := au.adminRepo.ViewCategory(category)

	if err != nil {
		return err
	} else {

		err := au.adminRepo.DeleteCategory(category)
		if err != nil {
			return err
		}
		return nil
	}
}

// Brands Management

func (au *adminUseCase) AddBrand(NewBrand domain.Brand) error {

	err := au.adminRepo.AddBrand(NewBrand)
	if err != nil {
		return err
	}
	return nil
}
func (au *adminUseCase) ListBrands() ([]domain.Brand, error) {

	brands, err := au.adminRepo.ListBrands()
	if err != nil {
		return nil, err
	}
	return brands, nil

}
func (au *adminUseCase) EditBrand(NewBrand domain.Brand) error {

	brand, err := au.adminRepo.ViewBrand(NewBrand.Brand_ID)

	if err != nil {
		return err
	} else {
		if brand.Brand_Name != NewBrand.Brand_Name {
			if err := au.adminRepo.EditBrand(NewBrand); err != nil {
				return err
			}
			return nil
		}
		return errors.New("entered brand name is same as the existing")
	}
}
func (au *adminUseCase) DeleteBrand(brand domain.Brand) error {

	brand, err := au.adminRepo.ViewBrand(brand.Brand_ID)

	if err != nil {
		return err
	} else {
		err := au.adminRepo.DeleteBrand(brand)
		if err != nil {
			return err
		}
		return nil
	}
}

// Product Management

func (au *adminUseCase) AddProducts(newProduct domain.Products) (string, error) {
	newProduct.Status = "available"
	product_code, err := au.adminRepo.AddProducts(newProduct)
	if err != nil {
		return "", err
	}

	return product_code, nil
}
func (au *adminUseCase) ListProducts(page, perPage int) ([]domain.ProductResponse, utils.MetaData, error) {

	products, metaData, err := au.adminRepo.ListProducts(page, perPage)
	if err != nil {
		return nil, metaData, err
	}

	return products, metaData, err
}
func (au *adminUseCase) EditProducts(newProduct domain.Products) error {

	err := au.adminRepo.EditProducts(newProduct)
	if err != nil {
		return err
	}

	return nil
}
func (au *adminUseCase) DeleteProducts(product_id string) error {

	err := au.adminRepo.DeleteProducts(product_id)
	if err != nil {
		return err
	}

	return nil
}

//Coupon

func (au *adminUseCase) AddCoupon(coupon domain.Coupon) error {
	err := au.adminRepo.AddCoupon(coupon)
	if err != nil {
		return err
	}

	return nil
}
func (au *adminUseCase) ListCoupon(page, perPage int) ([]domain.CouponResponse, utils.MetaData, error) {
	coupons, metaData, err := au.adminRepo.ListCoupon(page, perPage)
	if err != nil {
		return coupons, metaData, err
	}
	return coupons, metaData, nil
}
func (au *adminUseCase) DeleteCoupon(coupon_id string) error {
	err := au.adminRepo.DeleteCoupon(coupon_id)
	if err != nil {
		return err
	}

	return nil
}
