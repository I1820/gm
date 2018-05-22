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
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/configor"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"

	mgo "gopkg.in/mgo.v2"
)

// Config represents main configuration
var Config = struct {
	DB struct {
		URL string `default:"mongodb://172.18.0.1:27017" env:"db_url"`
	}
	Broker struct {
		URL string `default:"172.23.132.37:1884" env:"broker_url"`
	}
}{}

func main() {
	fmt.Println("GM AIoTRC @ 2018")

	// Load configuration
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}

	// Create a Mongo Session
	session, err := mgo.Dial(Config.DB.URL)
	if err != nil {
		log.Fatalf("Mongo session %s: %v", Config.DB.URL, err)
	}
	defer session.Close()
	fmt.Printf("Mongo session %s has been created\n", Config.DB.URL)

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

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
		ClientID: []byte(fmt.Sprintf("isrc-uplink-%d", rand.Int63())),
	}); err != nil {
		log.Fatalf("MQTT session %s: %s", Config.Broker.URL, err)
	}
	fmt.Printf("MQTT session %s has been created\n", Config.Broker.URL)

	// Subscribe to topics
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				// https://docs.loraserver.io/use/getting-started/
				TopicFilter: []byte("gateway/+/rx"),
				QoS:         mqtt.QoS0,
				Handler: func(topicName, message []byte) {
					fmt.Println(string(message))

					var m lora.GatewayMessage
					json.Unmarshal(message, &m)

					fmt.Println(m.PhyPayload.MHDR.MType)
					fmt.Println(m.PhyPayload.MHDR.RFU)
					fmt.Println(m.PhyPayload.MHDR.Major)

					fmt.Println(m.PhyPayload.FHDR.DevAddr)
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("MQTT subscription: %s", err)
	}

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Wait for receiving a signal.
	<-sigc

	fmt.Println("18.20 As always ... left me alone")
}
