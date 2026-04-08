package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emmanuella-codes/sms-mock/handlers"
	"github.com/emmanuella-codes/sms-mock/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	gin.SetMode(gin.ReleaseMode)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "http://localhost:4006"
	}

	s := store.New()
	h := handlers.New(s, webhookSecret, apiBaseURL)

	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/", h.UI)

	r.POST("/api/sms/send", h.Send)
	r.POST("/api/sms/inbound", h.Inbound)

	r.GET("/messages", h.AllMessages)
	r.GET("/messages/:phone", h.ByPhone)
	r.GET("/messages/:phone/latest", h.LatestMessage)
	r.GET("/otp/:phone", h.LatestOTP)

	r.DELETE("/messages", h.Clear)
	r.GET("/stats", h.Stats)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutCancel()
		if err := srv.Shutdown(shutCtx); err != nil {
			fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
		}
	}()
	fmt.Printf("\n  OLU SMS mock running\n")
	fmt.Printf("  UI       → http://localhost:%s\n", port)
	fmt.Printf("  Send     → POST http://localhost:%s/api/sms/send\n", port)
	fmt.Printf("  Inbound  → POST http://localhost:%s/api/sms/inbound\n", port)
	fmt.Printf("  Read OTP → GET  http://localhost:%s/otp/:phone\n", port)
	fmt.Printf("  All msgs → GET  http://localhost:%s/messages\n", port)
	fmt.Printf("  Latest   → GET  http://localhost:%s/messages/:phone/latest\n\n", port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
