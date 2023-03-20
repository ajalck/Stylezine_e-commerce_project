package usecase

import (
	"ajalck/e_commerce/domain"
	repoInt "ajalck/e_commerce/repository/interface"
	services "ajalck/e_commerce/usecase/interface"
	"errors"
)

type userAuthService struct {
	userRepo repoInt.UserRepository
}
type adminAuthService struct {
	adminRepo repoInt.AdminRepository
}

func NewUserAuthService(userRepo repoInt.UserRepository) services.UserAuth {
	return &userAuthService{userRepo: userRepo}
}

func NewAdminAuthService(adminRepo repoInt.AdminRepository) services.AdminAuth {
	return &adminAuthService{adminRepo: adminRepo}
}

func (ua *userAuthService) FindUser(email, role string) (domain.User, error) {
	user, err := ua.userRepo.FindUser(email, role)
	return user, err
}
func (ua *userAuthService) VerifyUser(email string, password string, userRole string) (bool, error) {

	user, err := ua.userRepo.FindUser(email, userRole)
	if err != nil {
		return false, err
	}
	// 	fmt.Println("\n\n", user.Password, "\n\n", password, "\n\n")
	// 	isValidPassword := CheckPasswordHash(user.Password, password)
	// 	fmt.(isValidPassword)
	// 	if !isValidPassword {
	// 		fmt.Println(isValidPassword)
	// 		return errors.New("failed to login ! Check your credentials")

	// 	}

	// 	return nil
	// }
	// func CheckPasswordHash(hash, password string) bool {
	// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// 	return err == nil

	hashedPassword := HashPassword(password)
	if user.Password != hashedPassword {
		return false, errors.New("Password not verified")
	}
	return true, nil
}
func (au *adminAuthService) FindAdmin(email, role string) (domain.User, error) {
	admin, err := au.adminRepo.FindAdmin(email, role)
	return admin, err
}
func (aa *adminAuthService) VerifyAdmin(email, password, userRole string) (bool, error) {

	user, err := aa.adminRepo.FindAdmin(email, userRole)
	if err != nil {
		return false, err
	}

	hashedPassword := HashPassword(password)
	if user.Password != hashedPassword {
		return false, errors.New("Missmatch in passwords")
	}
	return true, nil
}
