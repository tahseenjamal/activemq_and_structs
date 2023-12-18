package brokers

import (
	"fmt"
	"log"
	"time"

	"github.com/go-stomp/stomp"
)

// MessageBroker represents a simple message broker using the STOMP protocol.
type MessageBroker struct {
	conn           *stomp.Conn
	brokerURL      string
	username       string
	password       string
	heartbeat      time.Duration
	heartbeatGrace time.Duration
}

type MessageBrokerFunctions interface {
	Connect() error
	Disconnect() error
	Send(string, string) error
	Subscribe(string, func(*stomp.Message) error
}

// NewMessageBroker creates a new instance of MessageBroker.
func NewMessageBroker(brokerURL string, username string, password string, heartbeat, heartbeatGrace time.Duration) *MessageBroker {

	return &MessageBroker{
		conn:           nil,
		brokerURL:      brokerURL,
		username:       username,
		password:       password,
		heartbeat:      heartbeat,
		heartbeatGrace: heartbeatGrace,
	}
}

// Connect connects to the message broker.
func (mb *MessageBroker) Connect() error {
	if mb.conn != nil {
		return fmt.Errorf("already connected")
	}
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(mb.username, mb.password),
		stomp.ConnOpt.HeartBeat(mb.heartbeat, mb.heartbeatGrace),
	}

	conn, err := stomp.Dial("tcp", mb.brokerURL, options...)
	if err != nil {

		return err
	}

	mb.conn = conn
	return nil
}

// Disconnect disconnects from the message broker.
func (mb *MessageBroker) Disconnect() error {
	if mb.conn == nil {
		return fmt.Errorf("not connected")
	}

	fmt.Println("Disconnecting...")
	err := mb.conn.Disconnect()
	if err != nil {

		mb.conn = nil
	}
	return nil
}

// Send sends a message to a specified destination.
func (mb *MessageBroker) Send(destination, body string) error {
	if mb.conn == nil {
		return fmt.Errorf("not connected")
	}

	err := mb.conn.Send(destination, "text/plain", []byte(body))
	if err != nil {
		return fmt.Errorf("cannot send to queue: %s ", destination)
	}
	return nil
}

// Subscribe subscribes to messages from a specified destination.
func (mb *MessageBroker) Subscribe(destination string, callback func(*stomp.Message)) error {

	if mb.conn == nil {
		return fmt.Errorf("not connected to queue: %s ", destination)
	}

	sub, err := mb.conn.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Read()
			if err != nil {
				log.Fatalf("error reading from : %s, %s", destination, err.Error())
			}
			callback(msg)
		}
	}()

	return nil

}
