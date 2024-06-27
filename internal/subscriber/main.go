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
	rdb        *redis.Client
	subscriber *redis.PubSub
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
	router.GET("/health", healthHandler)

	srv := &http.Server{
		Addr:    ":3002",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	subscriber = rdb.Subscribe(context.Background(), channel)
	go func() {
		for msg := range subscriber.Channel() {
			fmt.Printf("Received message: %s\n", msg.Payload)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := subscriber.Close(); err != nil {
		log.Printf("Failed to unsubscribe from Redis channel: %v", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func healthHandler(c *gin.Context) {
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Redis connection failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Redis connection successful"})
}
