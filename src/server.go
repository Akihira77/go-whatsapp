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

	chatRepository := repositories.NewChatRepository(store)
	chatService := services.NewChatService(chatRepository)

	pageRouter(router, authenticatedPage, userService, chatService)
	userRouter(api, authenticatedApi, userService, chatService)
	chatRouter(authenticatedApi, chatService)

	hub := handlers.NewHub()
	go hub.Run()
	authenticatedApi.GET("/ws", func(c *gin.Context) {
		handlers.ServeWs(c, hub, userService, chatService)
	})
}

func pageRouter(router *gin.Engine, authenticatedPage *gin.RouterGroup, userService *services.UserService, chatService *services.ChatService) {
	pageHandler := handlers.NewPageHandler(userService, chatService)
	authenticatedPage.GET("/", pageHandler.RenderHome)
	authenticatedPage.GET("/chat", pageHandler.RenderChatPage)
	authenticatedPage.GET("/contacts", pageHandler.RenderMyContacts)
	authenticatedPage.GET("/users", pageHandler.RenderUsers)
	authenticatedPage.GET("/groups", pageHandler.RenderMakeGroup)
	authenticatedPage.GET("/groups/:groupId", pageHandler.RenderGroupInfo)
	authenticatedPage.GET("/groups/naming", pageHandler.RenderNamingGroup)
	authenticatedPage.GET("/users/profile", pageHandler.RenderMyProfile)
	authenticatedPage.GET("/users/edit", pageHandler.RenderEditProfile)
	authenticatedPage.GET("/users/change-password", pageHandler.RenderChangePassword)
	authenticatedPage.POST("/groups", pageHandler.CreateGroup)

	router.GET("/signup", pageHandler.RenderSignup)
	router.GET("/signin", pageHandler.RenderSignin)
}

func userRouter(api *gin.RouterGroup, authenticatedApi *gin.RouterGroup, userService *services.UserService, chatService *services.ChatService) {
	userHandler := handlers.NewUserHandler(userService, chatService)
	api.POST("/users/signin", userHandler.Signin)
	api.POST("/users/signup", userHandler.Signup)

	authenticatedApi.GET("/users", userHandler.GetUsers)
	authenticatedApi.GET("/users/my-image", userHandler.GetMyImageProfile)
	authenticatedApi.GET("/users/images/:userId", userHandler.GetUserImageProfile)
	authenticatedApi.GET("/groups/images/:groupId", userHandler.GetGroupImageProfile)
	authenticatedApi.GET("/users/my-info", userHandler.GetMyInfo)
	authenticatedApi.GET("/users/contacts", userHandler.GetMyContacts)
	authenticatedApi.GET("/groups/:groupId/members", userHandler.GetGroupMembers)
	authenticatedApi.PATCH("/users", userHandler.UpdateUserProfile)
	authenticatedApi.PATCH("/users/change-password", userHandler.UpdatePassword)
	authenticatedApi.POST("/users/contacts/:userId", userHandler.AddContact)
	authenticatedApi.POST("/users/logout", userHandler.Logout)
	authenticatedApi.DELETE("/users/contacts/:userId", userHandler.RemoveContact)
	authenticatedApi.PATCH("/groups/:groupId", userHandler.EditGroup)
}

func chatRouter(authenticatedApi *gin.RouterGroup, chatService *services.ChatService) {
	chatHandler := handlers.NewChatHandler(chatService)

	authenticatedApi.GET("/messages", chatHandler.GetChatList)
	authenticatedApi.GET("/messages/last/:username", chatHandler.SearchLastMessage)
	// authenticatedApi.POST("/messages", chatHandler.SendMsgToOffUser)
}
