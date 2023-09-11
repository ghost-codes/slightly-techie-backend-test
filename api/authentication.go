package api

import (
	"database/sql"
	"fmt"
	models "ghost-codes/slightly-techie-blog/db/models"
	"ghost-codes/slightly-techie-blog/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createUserWithEmailPasswordReq struct {
	Username  string `json:"username" binding:"required,min=8,alphanum"`
	FirstName string `json:"first_name" binding:"required,alphanum"`
	LastName  string `json:"last_name" binding:"required,alphanum"`
	Email     string `json:"email" binding:"required,email"`
	Contact   string `json:"contact"`
	Password  string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type loginUserResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

func newUserResponse(user models.User) UserResponse {
	return UserResponse{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
	}
}

// @Summary      signup new user
// @Description  create new user using email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body    body   	createUserWithEmailPasswordReq  true " "
// @Success      200  	{object}   	loginUserResponse
// @response     default  {object}  errorJson
// @Router       /signup [post]
func (server *Server) createUserWithEmailPassword(ctx *gin.Context) {
	req := createUserWithEmailPasswordReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrorJson(err))
		return
	}

	hasedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}

	createdUser := models.User{
		HashedPassword: hasedPassword,
		Username:       req.Username,
		Email:          req.Email,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
	}

	err = server.db.CreateUser(&createdUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}

	res := newUserResponse(createdUser)

	accessToken, _, err := server.tokenMaker.CreateToken(createdUser.ID, createdUser.SecurityKey, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return

	}
	refreshToken, _, err := server.tokenMaker.CreateToken(createdUser.ID, createdUser.SecurityKey, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return

	}

	ctx.JSON(http.StatusOK, loginUserResponse{User: res, RefreshToken: refreshToken, AccessToken: accessToken})
}

type loginUserReq struct {
	UsernameEmail string `json:"username_email" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

// @Summary      log existing user in
// @Description  log existing users in with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body    body   	loginUserReq  true " "
// @Success      200  	{object}   	loginUserResponse
// @response     default  {object}  	errorJson
// @Router       /login [post]
func (server *Server) loginWithEmailPassword(ctx *gin.Context) {
	req := loginUserReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrorJson(err))
		return
	}

	user, err := server.db.GetUserByEmail(req.UsernameEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("user does not exist:%v", err)
			ctx.JSON(http.StatusNotFound, NewErrorJson(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}

	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		err := fmt.Errorf("invalid credentials")
		ctx.JSON(http.StatusUnauthorized, NewErrorJson(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(user.ID, user.SecurityKey, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}
	refresh_token, _, err := server.tokenMaker.CreateToken(user.ID, user.SecurityKey, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}

	userRes := newUserResponse(*user)

	ctx.JSON(http.StatusOK, loginUserResponse{RefreshToken: refresh_token, AccessToken: accessToken, User: userRes})
}

type verifyEmailReq struct {
	ID   int64  `form:"id" binding:"required"`
	Code string `form:"code" binding:"required,min=64"`
}

// // @Summary      verify email
// // @Description  verify email using short code
// // @Tags         Authentication
// // @Accept       json
// // @Produce      json
// // @Param        body    body   	loginUserReq  true " "
// // @Success      200  	{object}    UserResponse
// // @response     default  {object}  	errorJson
// // @Router       /verify_email [post]
// func (server *Server) verifyEmail(ctx *gin.Context) {
// 	var req verifyEmailReq

// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, NewErrorJson(err))
// 		return
// 	}

// 	verifyEmail, err := server.db.GetVerifyEmail(ctx, req.ID)
// 	if err != nil || verifyEmail.SecretKey != req.Code {
// 		err := fmt.Errorf("invalid verification link:%w", err)
// 		ctx.JSON(http.StatusBadRequest, NewErrorJson(err))
// 		return
// 	}

// 	if time.Now().After(verifyEmail.ExpiredAt) {
// 		err := fmt.Errorf("verification link has expired")
// 		ctx.JSON(http.StatusBadRequest, NewErrorJson(err))
// 		return
// 	}

// 	// replace user's email when user data email is different from verify_email email
// 	user, err := server.db.GetUser(ctx, *verifyEmail.Username)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
// 		return
// 	}

// 	now := time.Now()
// 	arg := db.UpdateUserParams{
// 		ID:                user.ID,
// 		Username:          user.Username,
// 		Email:             verifyEmail.Email,
// 		FirstName:         user.FirstName,
// 		LastName:          user.LastName,
// 		HashedPassword:    user.HashedPassword,
// 		AvatarUrl:         user.AvatarUrl,
// 		Contact:           user.Contact,
// 		SecurityKey:       uuid.NewString(),
// 		VerifiedAt:        &now,
// 		TwitterSocial:     user.TwitterSocial,
// 		AppleSocial:       user.AppleSocial,
// 		GoogleSocial:      user.GoogleSocial,
// 		CreatedAt:         user.CreatedAt,
// 		PasswordChangedAt: user.PasswordChangedAt,
// 	}

// 	updatedUser, err := server.db.UpdateUser(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, newUserResponse(updatedUser))
// }

// // @Summary      send verification email
// // @Description  short code is sent to the user's email for verification
// // @Tags         Authentication
// // @Security 	bearerAuth
// // @Accept       json
// // @Produce      json
// // @Success      200  	string    "verification email successfully sent"
// // @response     default  {object}  	errorJson
// // @Router       /send_verification_email [get]
// func (server *Server) sendVerificationEmail(ctx *gin.Context) {
// 	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

// 	err := server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, &worker.PayloadSendVerifyEmail{
// 		UserID: payload.UserID,
// 	})

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, map[string]string{"success": "verification email successfully sent"})
// }
