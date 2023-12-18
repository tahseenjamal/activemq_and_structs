// main.go
package main

import (
	"fmt"
	"log"
	"messagebroker/brokers"
	"time"

	"github.com/go-stomp/stomp"
)

var (
	messageBroker *brokers.MessageBroker
)

func main() {

	// If you are running dockers in mac, the external host
	// localhost is referred to as host.docker.internal
	brokerURL := "host.docker.internal:61613"

	// In case of linux, you can use below
	// brokerURL := "localhost:61613"
	username := "admin"
	password := "admin"
	messageBroker = brokers.NewMessageBroker(brokerURL, username, password, 5*time.Second, 5*time.Second)

	err := messageBroker.Connect()
	if err != nil {
		log.Fatal("Error connecting to message broker:", err)
	}

	defer func() {
		err := messageBroker.Disconnect()
		if err != nil {
			log.Fatal("Error disconnecting from message broker:", err)
		}
	}()

	// Subscribe to a destination
	err = messageBroker.Subscribe("/example", func(msg *stomp.Message) {
		fmt.Println("Received message:", string(msg.Body))
	})
	if err != nil {
		log.Fatal("Error subscribing to destination:", err)
	}

	for {

		// Send a message
		err = messageBroker.Send("/example", "Again Hello, Message Broker!")
		if err != nil {
			log.Fatal("Error sending message:", err)
		}
		time.Sleep(1 * time.Second)
	}
}
