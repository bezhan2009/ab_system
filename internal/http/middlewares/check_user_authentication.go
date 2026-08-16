package middlewares

import (
	"ab_system/internal/http"
	"ab_system/internal/http/response"
	"ab_system/internal/lib/jwt"
	"ab_system/pkg/errs"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckUserAuthentication(c *gin.Context) {
	header := c.GetHeader(http.AuthorizationHeader)

	if header == "" {
		response.HandleError(c, errs.ErrInvalidToken)
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {

		response.HandleError(c, errs.ErrInvalidToken)
		return
	}

	if len(headerParts[1]) == 0 {
		response.HandleError(c, errs.ErrInvalidToken)
		return
	}

	accessToken := headerParts[1]

	claims, err := jwtauth.ParseToken(accessToken)
	if err != nil {
		response.HandleError(c, errs.ErrInvalidToken)
		return
	}

	c.Set(http.UserIDCtx, claims.Sub)
	c.Set(http.UserRoleCtx, claims.Role)

	c.Next()
}
