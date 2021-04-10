package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"server/middleware/gin"
	"syscall"
	"time"
)

//
// Public
//
func Run() {
	// port
	port := getPort()

	// Address to bind to
	bind := getBindAddr()

	// Setup gin handler
	handler := gin.Setup()

	// HTTP Server
	server := &http.Server{
		Addr:           bind + ":" + port,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Handle graceful shutdown on SIGINT
	idleConnsClosed := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt, syscall.SIGTERM)
		<-s

		fmt.Printf("Http server is shutdowning")

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			fmt.Printf("HTTP server shutdown error: %v\n", err)
		}
		close(idleConnsClosed)
	}()

	// Start the server
	fmt.Printf("Starting server on http://%s\n", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}

	<-idleConnsClosed
}

//
//  private
//

func getPort() string {
	if len(os.Getenv("PORT")) != 0 {
		return os.Getenv("PORT")
	}
	return "8080" // default
}

func getBindAddr() string {
	if len(os.Getenv("BIND")) != 0 {
		return os.Getenv("BIND")
	}
	return "" // default
}
