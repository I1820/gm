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
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"

	"github.com/aiotrc/gm/lora"
	"github.com/brocaar/lorawan"
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/configor"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Config represents main configuration
var Config = struct {
	Broker struct {
		URL string `default:"172.23.132.37:1884" env:"broker_url"`
	}
	Device struct {
		Addr    string `default:"00000030"`
		AppSKey [16]byte
		NetSKey [16]byte
	}
}{}

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
					fmt.Println(string(message))
					fmt.Println()

					var m lora.GatewayMessage
					if err := json.Unmarshal(message, &m); err != nil {
						log.Error(err)
						return
					}

					macPayload := m.PhyPayload.MACPayload.(*lorawan.MACPayload)
					if Config.Device.Addr == fmt.Sprintf("%v", macPayload.FHDR.DevAddr) {
						ok, err := m.PhyPayload.ValidateMIC(Config.Device.NetSKey)
						if err != nil {
							log.Error(err)
						}
						if !ok {
							log.Error("Invalid MIC")
						}
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
