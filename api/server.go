package api

import (
	"fmt"

	config "ghost-codes/slightly-techie-blog/config"
	db "ghost-codes/slightly-techie-blog/db"
	middleware "ghost-codes/slightly-techie-blog/middleware"
	token "ghost-codes/slightly-techie-blog/token"

	models "ghost-codes/slightly-techie-blog/db/models"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/stripe/stripe-go/v72"
)

type Server struct {
	config     config.Config
	router     *gin.Engine
	tokenMaker token.Maker
	db         *models.Store
}

func NewServer(config config.Config) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.SecretKey)

	stripe.Key = config.SecretKey

	if err != nil {
		return nil, err
	}
	db, err := db.NewGorm(config.DBSource())
	if err != nil {
		return nil, err
	}

	store := models.Store{
		DB: db,
	}

	server := &Server{
		config: config,

		tokenMaker: maker,
		db:         &store,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("phone", validPhoneNumber)
	// }

	router.Use(middleware.Errors("./config/.env", "rollbackToken", service.ConsoleLogger))

	user := router.Group("/user")

	{

		user.POST("/sign_up", server.createUserWithEmailPassword)
		user.POST("/login", server.loginWithEmailPassword)
		auth := user.Use(middleware.AuthMiddleware(*server.db, server.tokenMaker))
		{
			auth.GET("/post", server.ViewAllPosts)
			auth.POST("/post", server.CreatePost)
			auth.GET("/post/:id", server.ViewPost)
			auth.DELETE("/post/:id", server.DeletePost)
			auth.PATCH("/post/:id", server.UpdatePost)

		}

	}

	// driver := router.Group("/driver")
	// {
	// 	driverRepo := NewDriverRepo(server)
	// 	driver.POST("/create", driverRepo.create)
	// 	driver.POST("/login", driverRepo.login)
	// 	auth := driver.Group("/", middleware.DriverAuthMiddlware(*server.db, server.tokenMaker))
	// 	{
	// 		auth.GET("me", func(ctx *gin.Context) {
	// 			driver := ctx.MustGet(middleware.DriverPayloadKey).(models.Driver)

	// 			ctx.JSON(http.StatusOK, newDriverResponse(driver))

	// 		})
	// 	}
	// }

	server.router = router
}

func (server *Server) Start(addr string) error {
	fmt.Println("======= Server has Started ==============")
	return server.router.Run(addr)

}

type errorJson struct {
	Message string `json:"message"`
}

func NewErrorJson(err error) gin.H {
	return gin.H{
		"message": err.Error(),
	}
}

func GenericSuccesMsg(msg string) gin.H {
	return gin.H{
		"message": msg,
	}
}
