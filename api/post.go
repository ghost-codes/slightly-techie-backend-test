package api

import (
	models "ghost-codes/slightly-techie-blog/db/models"
	"ghost-codes/slightly-techie-blog/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createPostRequest struct {
	Text string `json:"text" binding:"min=10,required"`
}

type updatePostRequest struct {
	Text string `json:"text" binding:"min=10,required"`
}
type idParam struct {
	ID int64 `uri:"id" binding:"required"`
}

type postResponce struct {
	Data    models.Post `json:"data"`
	Message string      `json:"message"`
}
type postsResponse struct {
	Data    []models.Post `json:"data"`
	Message string        `json:"message"`
}

func (server *Server) CreatePost(ctx *gin.Context) {
	user, err := util.GetUserFromCTX(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, NewErrorJson(err))
		return
	}
	req := createPostRequest{}

	if err := ctx.BindJSON(&req); err != nil {
		return
	}

	post := models.Post{
		Text:   req.Text,
		UserID: user.ID,
	}

	if err := server.db.CreatePost(&post); err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}

	ctx.JSON(http.StatusOK, postResponce{
		Data:    post,
		Message: "Post created successfully",
	})
}

func (server *Server) UpdatePost(ctx *gin.Context) {
	user, err := util.GetUserFromCTX(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, NewErrorJson(err))
		return
	}
	req := updatePostRequest{}
	param := idParam{}

	if err := ctx.BindJSON(&req); err != nil {
		return
	}
	if err := ctx.BindUri(&param); err != nil {
		return
	}

	post := models.Post{
		Text:   req.Text,
		UserID: user.ID,
		ID:     param.ID,
	}

	if err := server.db.UpdatePost(&post); err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}

	ctx.JSON(http.StatusOK, postResponce{
		Data:    post,
		Message: "Post updated successfully",
	})
}

func (server *Server) DeletePost(ctx *gin.Context) {
	user, err := util.GetUserFromCTX(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, NewErrorJson(err))
		return
	}
	req := idParam{}

	if err := ctx.BindUri(&req); err != nil {
		return
	}

	if err := server.db.DeletePost(req.ID, user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}

func (server *Server) ViewPost(ctx *gin.Context) {
	user, err := util.GetUserFromCTX(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, NewErrorJson(err))
		return
	}
	req := idParam{}

	if err := ctx.BindUri(&req); err != nil {
		return
	}

	post, err := server.db.ViewOnePost(req.ID, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}

	ctx.JSON(http.StatusOK, postResponce{
		Data:    *post,
		Message: "success",
	})
}

func (server *Server) ViewAllPosts(ctx *gin.Context) {
	user, err := util.GetUserFromCTX(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, NewErrorJson(err))
		return
	}

	post, err := server.db.ViewAll(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}

	ctx.JSON(http.StatusOK, postsResponse{
		Data:    post,
		Message: "success",
	})
}
