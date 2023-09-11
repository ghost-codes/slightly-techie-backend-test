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

type queryRequestParams struct {
	Count int32 `form:"count" binding:"default=20,min=5"`
	Page  int32 `form:"page" binding:"default=1"`
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

// @Summary      Create Post
// @Description  create new user post
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body    body   	createPostRequest  true " "
// @Success      200  	{object}   	postResponce
// @response     default  {object}  errorJson
// @Router       /post [post]
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

// @Summary      Update Post
// @Description  update user post
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body    body   	updatePostRequest  true " "
// @Param      	id 		path	int true " "
// @Success      200  	{object}   	postResponce
// @response     default  {object}  errorJson
// @Router       /post/{id} [patch]
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

// @Summary      Delete Post
// @Description  delete user post
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param      	id 		path	int true " "
// @response     default  {object}  errorJson
// @Router       /post/{id} [delete]
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

// @Summary      View Post
// @Description  view user post with id
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param      	id 		path	int true " "
// @Param      count 		query	int false " "
// @Param      	page 		query	int false " "
// @Success      200  	{object}   	postsResponse
// @response     default  {object}  errorJson
// @Router       /post/{id} [get]
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

// @Summary      View Post
// @Description  view user post with id
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @response     default  {object}  errorJson
// @Router       /post [get]
func (server *Server) ViewAllPosts(ctx *gin.Context) {
	queryParams := queryRequestParams{}

	if err := ctx.BindQuery(&queryParams); err != nil {
		return
	}
	user, err := util.GetUserFromCTX(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, NewErrorJson(err))
		return
	}

	post, err := server.db.ViewAll(user.ID, int(queryParams.Count), int(queryParams.Page*queryParams.Count))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorJson(err))
		return
	}

	ctx.JSON(http.StatusOK, postsResponse{
		Data:    post,
		Message: "success",
	})
}
