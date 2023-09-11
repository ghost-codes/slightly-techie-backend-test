package middleware

import (
	"fmt"
	"net/http"
	"strings"

	db "ghost-codes/slightly-techie-blog/db/models"
	token "ghost-codes/slightly-techie-blog/token"

	"github.com/gin-gonic/gin"
)

const (
	authorizationKey               = "Authorization"
	authorizationBearerType        = "bearer"
	UserPayloadKey                 = "user"
	DriverPayloadKey               = "driver"
	UserTypeClient          string = "client"
	UserTypeKey                    = "userType"
)

func AuthMiddleware(store db.Store, maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr, err := extractTokenString(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newErrorJson(err))
			return
		}

		payload, err := maker.VerifyToken(*tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newErrorJson(err))
			return
		}

		//TODO:
		user, err := store.GetUserByID(payload.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, newErrorJson(fmt.Errorf("Unathenticated")))
			return
		}

		ctx.Set(UserPayloadKey, user)

		ctx.Next()
	}
}

func extractTokenString(ctx *gin.Context) (*string, error) {
	authorizationHeader := ctx.GetHeader(authorizationKey)
	if len(authorizationHeader) == 0 {

		return nil, fmt.Errorf("Unathenticated")
	}
	fields := strings.Fields(authorizationHeader)

	if len(fields) < 2 {

		return nil, fmt.Errorf("Authorization bearer type mismatch")
	}

	if strings.ToLower(fields[0]) != authorizationBearerType {

		return nil, fmt.Errorf("Authorization bearer type mismatch")
	}

	return &fields[1], nil
}

func newErrorJson(err error) gin.H {
	return gin.H{
		"message": err.Error(),
	}
}
