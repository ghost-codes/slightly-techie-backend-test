package util

import (
	"errors"

	db "ghost-codes/slightly-techie-blog/db/models"
	middleware "ghost-codes/slightly-techie-blog/middleware"

	"github.com/gin-gonic/gin"
)

func GetUserFromCTX(ctx *gin.Context) (user *db.User, err error) {
	u := ctx.MustGet(middleware.UserPayloadKey).(*db.User)
	user = u
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("Unauthorized")
		}
	}()
	return
}
