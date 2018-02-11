/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-11-2017
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aiotrc/gm/gateway"
	"github.com/gin-gonic/gin"
)

var gateways []gateway.Gateway

func init() {
	gateways = make([]gateway.Gateway, 0)
}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)

		api.POST("/gateway", gatewayNewHandler)
		api.GET("/gateway", gatewayListHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	return r
}

func main() {
	fmt.Println("GM AIoTRC @ 2018")

	r := handle()

	srv := &http.Server{
		Addr:    ":1373",
		Handler: r,
	}

	go func() {
		fmt.Printf("GM Listen: %s\n", srv.Addr)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Listen Error:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("GM Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown Error:", err)
	}
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")
}

func gatewayListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gateways)
}

func gatewayNewHandler(c *gin.Context) {
	var json gatewayReq
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := json.Name
	address := json.Address

	g, err := gateway.New(name, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gateways = append(gateways, g)

	c.JSON(http.StatusOK, g)
}
