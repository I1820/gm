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
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brocaar/lorawan"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "404 Not Found"})
	})

	r.Use(gin.ErrorLogger())

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)
		api.POST("/decrypt", decryptHandler)
	}

	return r
}

func main() {
	fmt.Println("GM AIoTRC @ 2018")

	srv := &http.Server{
		Addr:    ":1374",
		Handler: handle(),
	}

	go func() {
		fmt.Printf("GM Listen: %s\n", srv.Addr)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Listen Error:", err)
		}
	}()

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	// Wait for receiving a signal.
	<-sigc

	fmt.Println("18.20 As always ... left me alone")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown Error:", err)
	}
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")

}

func decryptHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var json decryptReq
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	appSKeySlice, err := hex.DecodeString(json.AppSKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var appSKey lorawan.AES128Key
	copy(appSKey[:], appSKeySlice)

	netSKeySlice, err := hex.DecodeString(json.NetSKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var netSKey lorawan.AES128Key
	copy(netSKey[:], netSKeySlice)

	var phy lorawan.PHYPayload
	if err := phy.UnmarshalBinary(json.PhyPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	mac, ok := phy.MACPayload.(*lorawan.MACPayload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "*MACPayload expected"})
		return
	}

	success, err := phy.ValidateUplinkJoinMIC(netSKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid MIC"})
		return
	}

	if err := phy.DecryptFRMPayload(appSKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	data, ok := mac.FRMPayload[0].(*lorawan.DataPayload)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "*DataPayload expected"})
		return
	}

	c.JSON(http.StatusOK, data.Bytes)
}
