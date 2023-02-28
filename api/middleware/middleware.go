package middleware

import (
	services "ajalck/e_commerce/usecase/interface"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserMiddleware interface {
	AuthorizeJWT(c *gin.Context)
}
type middleware struct {
	jwtService services.JwtServices
}

func NewUserMiddleware(jwtServices services.JwtServices) UserMiddleware {
	return &middleware{
		jwtService: jwtServices,
	}
}
func (m *middleware) AuthorizeJWT(c *gin.Context) {

	authHeader := c.Request.Header.Get("Authorization")
	bearerToken := strings.Split(authHeader, " ")

	if len(bearerToken) != 2 {
		c.JSON(400, "Something went wrong ! Please login first if you haven't")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}
	authToken := bearerToken[1]
	ok, claims := m.jwtService.VerifyToken(authToken)
	if !ok {
		c.JSON(400, "Autherization failed !")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}
	user_id := fmt.Sprint(claims.UserId)
	user_name := fmt.Sprint(claims.Username)
	user_role := fmt.Sprint(claims.UserRole)
	c.Writer.Header().Set("id", user_id)
	c.Writer.Header().Set("username", user_name)
	c.Writer.Header().Set("userrole", user_role)
}
