package main

import (
	"encoding/json"
	"github.com/applait/xplex-meta/contract"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func rmqPublish(rmqChan *amqp.Channel, exchange string, message map[string]string) {
	var m interface{}

	switch exchange {
	case "agent_register":
		hostname, _ := os.Hostname()
		m = contract.AgentRegister{
			Hostname: hostname,
			Secret:   viper.GetString("rmq.secret"),
		}

	case "agent_stream-init":
		m = contract.AgentStreamInit{
			Name:  message["StreamName"],
			Token: message["StreamToken"],
			JWT:   JWT,
		}

	case "default":
		m = json.RawMessage("")

	}

	publishMessage, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}

	// Publish messages to exchange
	publishErr := rmqChan.Publish(
		exchange, // Exchange
		"",       // Key
		false,    // Mandatory
		false,    // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(publishMessage),
		})

	if publishErr != nil {
		log.Println("Error publishing message to exchange"+exchange, publishErr)
	}
}
