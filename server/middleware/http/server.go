package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"server/config"
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

	// read Timeout
	readTimeout := getReadTimeout() * time.Second

	// write Timeout
	writeTimeout := getWriteTimeout() * time.Second

	// Setup gin handler
	handler := gin.Setup()

	// HTTP Server
	server := &http.Server{
		Addr:           bind + ":" + port,
		Handler:        handler,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
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
	if config.HasKey("http", "port") {
		val, _ := config.GetConfig().String("http", "port")
		return val
	}
	return "8080" // default
}

func getBindAddr() string {
	if len(os.Getenv("BIND")) != 0 {
		return os.Getenv("BIND")
	}

	if config.HasKey("http", "bind") {
		val, _ := config.GetConfig().String("http", "bind")
		return val
	}
	return "" // default
}

func getReadTimeout() time.Duration {
	if config.HasKey("http", "ReadTimeout") {
		val, _ := config.GetConfig().Int("http", "ReadTimeout")
		return time.Duration(val)
	}
	return 10 // default
}

func getWriteTimeout() time.Duration {
	if config.HasKey("http", "WriteTimeout") {
		val, _ := config.GetConfig().Int("http", "WriteTimeout")
		return time.Duration(val)
	}
	return 10 // default
}
