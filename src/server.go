package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Akihira77/go_whatsapp/src/store"
	"github.com/gin-gonic/gin"
)

func NewHttpServer(port string, store *store.Store) error {
	router := gin.Default()

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return srv.ListenAndServe()
}

func mainRouter(router *gin.Engine) {
	api := router.Group("/api")
	userRouter := api.Group("/user")

	userRouter.GET("/my-info")
	userRouter.POST("/sign-in")
	userRouter.GET("/refresh-token")
}
