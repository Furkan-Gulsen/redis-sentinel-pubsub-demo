package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	redisAddr = "haproxy:6379"
	channel   = "pubsub-channel"
)

var (
	rdb *redis.Client
)

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	fmt.Println("Redis connection successful")

	router := gin.Default()
	router.POST("/publish", publishHandler)
	router.GET("/health", healthHandler)

	srv := &http.Server{
		Addr:    ":3001",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

type publishRequest struct {
	Username string `json:"username" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

func publishHandler(c *gin.Context) {
	var req publishRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := fmt.Sprintf("%s: %s", req.Username, req.Message)
	err := rdb.Publish(context.Background(), channel, msg).Err()
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Failed to publish message"})
		return
	}

	rdb.Set(context.Background(), req.Username+time.Now().Format(time.RFC3339), msg, 1*time.Minute)
	fmt.Printf("Message sent: %s\n", msg)
	c.JSON(http.StatusOK, gin.H{"status": "Message sent", "message": msg})
}

func healthHandler(c *gin.Context) {
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Redis connection failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Redis connection successful"})
}
