package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func Serv(router *gin.Engine, addr string) {

	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    1800 * time.Second,
		WriteTimeout:   1800 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error: ", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("shutdown error: ", err.Error())
	}
}
