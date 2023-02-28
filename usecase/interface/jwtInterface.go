package interfaces

import "ajalck/e_commerce/domain"

type JwtServices interface {
	GenerateToken(user_id int, username, role string) string
	VerifyToken(token string) (bool, *domain.SignedDetails)
}
