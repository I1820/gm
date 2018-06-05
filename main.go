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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"

	"github.com/aiotrc/gm/lora"
	"github.com/brocaar/lorawan"
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/configor"

	"github.com/gin-gonic/gin"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Config represents main configuration
var Config = struct {
	Broker struct {
		URL string `default:"172.23.132.37:1884" env:"broker_url"`
	}
	Device struct {
		// Addr string `default:"2601146f"`
		Addr string `default:"0000003a"`
		// AppSKey [16]byte `default:"[0x29, 0xCB, 0xD0, 0x5A, 0x4C, 0xB9, 0xFB, 0xC5, 0x16, 0x6A, 0x89, 0xE6, 0x71, 0xC0, 0xEF, 0xCE]"`
		AppSKey [16]byte `default:"[0x2B, 0x7E, 0x15, 0x16, 0x28, 0xAE, 0xD2, 0xA6, 0xAB, 0xF7, 0x15, 0x88, 0x09, 0xCF, 0x4F, 0x3C]"`
		// NetSKey [16]byte `default:"[0x5E, 0xD4, 0x38, 0xE5, 0xC8, 0x6E, 0xDD, 0x00, 0xCE, 0x0E, 0xD6, 0x22, 0x2A, 0x99, 0xE6, 0x84]"`
		NetSKey [16]byte `default:"[0x2B, 0x7E, 0x15, 0x16, 0x28, 0xAE, 0xD2, 0xA6, 0xAB, 0xF7, 0x15, 0x88, 0x09, 0xCF, 0x4F, 0x3C]"`
	}
}{}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
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

	// Load configuration
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}

	// Create an MQTT client
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			log.WithFields(log.Fields{
				"component": "gm",
			}).Errorf("MQTT Client %s", err)
		},
	})
	defer cli.Terminate()

	// Connect to the MQTT Server.
	if err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  Config.Broker.URL,
		ClientID: []byte(fmt.Sprintf("isrcgm-%d", rand.Int63())),
	}); err != nil {
		log.Fatalf("MQTT session %s: %s", Config.Broker.URL, err)
	}
	fmt.Printf("MQTT session %s has been created\n", Config.Broker.URL)

	// Subscribe to topics
	if err := cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				// https://docs.loraserver.io/use/getting-started/
				TopicFilter: []byte("gateway/+/rx"),
				QoS:         mqtt.QoS0,
				Handler: func(topicName, message []byte) {
					log.Info(string(message))

					var m lora.GatewayMessage
					if err := json.Unmarshal(message, &m); err != nil {
						log.Error(err)
						return
					}

					phyJSON, err := m.PhyPayload.MarshalJSON()
					if err != nil {
						log.Error(err)
					}
					log.Info(string(phyJSON))

					macPayload, ok := m.PhyPayload.MACPayload.(*lorawan.MACPayload)
					if !ok {
						log.Error("*MACPayload expected")
					}

					log.Infof("DevAddr: %v", macPayload.FHDR.DevAddr)
					if Config.Device.Addr == fmt.Sprintf("%v", macPayload.FHDR.DevAddr) {
						ok, err := m.PhyPayload.ValidateMIC(Config.Device.NetSKey)
						if err != nil {
							log.Error(err)
						}
						if !ok {
							log.Error("Invalid MIC")
						}

						if err := m.PhyPayload.DecryptFRMPayload(Config.Device.AppSKey); err != nil {
							log.Error(err)
						}

						pl, ok := macPayload.FRMPayload[0].(*lorawan.DataPayload)
						if !ok {
							log.Error("*DataPayload expected")
							return
						}

						log.Info(pl.Bytes)
					}
				},
			},
		},
	}); err != nil {
		log.Fatalf("MQTT subscription: %s", err)
	}

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Wait for receiving a signal.
	<-sigc

	fmt.Println("18.20 As always ... left me alone")
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")

}

func decryptHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var json decryptReq
	if err := c.BindJSON(&json); err != nil {
		return
	}

	appSKeySlice, err := hex.DecodeString(json.AppSKey)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var appSKey lorawan.AES128Key
	copy(appSKey[:], appSKeySlice[:])

	netSKeySlice, err := hex.DecodeString(json.NetSKey)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var netSKey lorawan.AES128Key
	copy(netSKey[:], netSKeySlice[:])

	var phy lorawan.PHYPayload
	if err := phy.UnmarshalBinary(json.PhyPayload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	mac, ok := phy.MACPayload.(*lorawan.MACPayload)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("*MACPayload expected"))
		return
	}

	success, err := phy.ValidateMIC(netSKey)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !success {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Invalid MIC"))
		return
	}

	if err := phy.DecryptFRMPayload(appSKey); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data, ok := mac.FRMPayload[0].(*lorawan.DataPayload)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("*DataPayload expected"))
		return
	}

	c.JSON(http.StatusOK, data.Bytes)
}
