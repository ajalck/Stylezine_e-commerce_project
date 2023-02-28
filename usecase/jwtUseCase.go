package usecase

import (
	"ajalck/e_commerce/domain"
	services "ajalck/e_commerce/usecase/interface"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type jwtService struct {
	SecretKey string
}

func NewJWTService() services.JwtServices {
	return &jwtService{
		SecretKey: os.Getenv("SECRET_KEY"),
	}
}
func (j *jwtService) GenerateToken(userId int, email, role string) string {

	claims := &domain.SignedDetails{
		UserId:   userId,
		Username: email,
		UserRole: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))

	if err != nil {
		log.Println("\n\n\n", err, "Could'nt sign the token")
	}
	return signedToken
}

func (j *jwtService) VerifyToken(SignedToken string) (bool, *domain.SignedDetails) {
	claims := &domain.SignedDetails{}
	token, _ := j.GetTokenFromString(SignedToken, claims)
	if token.Valid {
		if t := claims.Valid(); t == nil {
			return true, claims
		}
	}
	_ = &domain.User{Verification: true}
	return false, claims
}
func (j *jwtService) GetTokenFromString(signedToken string, claims *domain.SignedDetails) (*jwt.Token, error) {
	return jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method:#{token.Header['alg']}")
		}
		return []byte(j.SecretKey), nil
	})
}
