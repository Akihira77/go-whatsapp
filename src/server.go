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
	router := gin.Default(func(e *gin.Engine) {
		e.MaxMultipartMemory = 2 << 20 // 20 MiB
	})

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
	router.Static("/public", "./public")

	userRepository := repositories.NewUserRepository(store)
	userService := services.NewUserService(userRepository)

	api := router.Group("/api")
	authenticatedApi := api.Group("", func(c *gin.Context) {
		middlewares.AuthOnly(c, userService)
	})
	authenticatedPage := router.Group("", func(c *gin.Context) {
		middlewares.AuthOnly(c, userService)
	})

	pageRouter(router, authenticatedPage, userService)
	userRouter(api, authenticatedApi, userService)
}

func pageRouter(router *gin.Engine, authenticatedPage *gin.RouterGroup, userService *services.UserService) {
	pageHandler := handlers.NewPageHandler(userService)
	authenticatedPage.GET("/", pageHandler.RenderHome)
	authenticatedPage.GET("/users/profile", pageHandler.RenderMyProfile)
	authenticatedPage.GET("/users/edit", pageHandler.RenderEditProfile)
	authenticatedPage.GET("/users/change-password", pageHandler.RenderChangePassword)

	router.GET("/signup", pageHandler.RenderSignup)
	router.GET("/signin", pageHandler.RenderSignin)
}

func userRouter(api *gin.RouterGroup, authenticatedApi *gin.RouterGroup, userService *services.UserService) {
	userHandler := handlers.NewUserHandler(userService)
	api.POST("/users/signin", userHandler.Signin)
	api.POST("/users/signup", userHandler.Signup)

	authenticatedApi.GET("/users", userHandler.GetUsers)
	authenticatedApi.GET("/users/my-image", userHandler.GetMyImageProfile)
	authenticatedApi.GET("/users/my-info", userHandler.GetMyInfo)
	authenticatedApi.GET("/users/contacts", userHandler.GetMyContacts)
	authenticatedApi.PATCH("/users", userHandler.UpdateUserProfile)
	authenticatedApi.PATCH("/users/change-password", userHandler.UpdatePassword)
	authenticatedApi.POST("/users/contacts/:userId", userHandler.AddContact)
	authenticatedApi.DELETE("/users/contacts/:userId", userHandler.RemoveContact)
}
