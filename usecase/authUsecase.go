package usecase

import (
	"ajalck/e_commerce/domain"
	repoInt "ajalck/e_commerce/repository/interface"
	services "ajalck/e_commerce/usecase/interface"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (ua *userAuthService) FindUser(c *gin.Context, email, role string) (domain.User, error) {
	user, err := ua.userRepo.FindUser(c, email, role)
	return user, err
}
func (ua *userAuthService) VerifyUser(c *gin.Context, email string, password string, userRole string) (bool, error) {

	user, err := ua.userRepo.FindUser(c, email, userRole)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found ! Please check your credentials"})
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
		c.JSON(400, "missmatch in passwords")
		return false, nil
	}
	// c.JSON(200, "Verified successfully")
	return true, nil
}
func (aa *adminAuthService) VerifyAdmin(c *gin.Context, email, password, userRole string) (bool, error) {

	user, err := aa.adminRepo.FindAdmin(c, email, userRole)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found ! Please check your credentials"})
		return false, err
	}

	hashedPassword := HashPassword(password)
	if user.Password != hashedPassword {
		c.JSON(400, "missmatch in passwords")
		return false, nil
	}
	c.JSON(200, "Verified successfully")
	return true, nil
}
