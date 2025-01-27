package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Akihira77/go_whatsapp/src/handlers"
	"github.com/Akihira77/go_whatsapp/src/middlewares"
	"github.com/Akihira77/go_whatsapp/src/repositories"
	"github.com/Akihira77/go_whatsapp/src/services"
	"github.com/Akihira77/go_whatsapp/src/store"
	"github.com/gin-gonic/gin"
)

func NewHttpServer(port string, store *store.Store) error {
	router := gin.Default()

	mainRouter(router, store)

	slog.Info("Run server", "port:", port)
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return srv.ListenAndServe()
}

func mainRouter(router *gin.Engine, store *store.Store) {
	userRepository := repositories.NewUserRepository(store)
	userService := services.NewUserService(userRepository)

	api := router.Group("/api")
	requireAuth := api.Group("", func(c *gin.Context) {
		middlewares.AuthOnly(c, userService)
	})

	userRouter(api, requireAuth, userService)
}

func userRouter(api *gin.RouterGroup, requireAuth *gin.RouterGroup, userService *services.UserService) {
	userHandler := handlers.NewUserHandler(userService)
	api.POST("/users/signin", userHandler.Signin)
	api.POST("/users/signup", userHandler.Signup)

	requireAuth.GET("/users/my-info", userHandler.GetMyInfo)
	requireAuth.GET("/users/contacts", userHandler.GetMyContacts)
}
