package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type RmqHandler struct {
	ExchangeRegister struct {
		Hostname string `json:"hostname"`
		Secret   string `json:"secret"`
	}

	ExchangeStreamInit struct {
		StreamName  string `json:"stream_name"`
		StreamToken string `json:"stream_token"`
		JWT         string `json:"jwt"`
	}
}

func (rmq RmqHandler) Send(rmqChan *amqp.Channel, exchange string, message map[string]string) {
	switch exchange {
	case "agent_register":
		response := q.ExchangeRegister{}

		hostname, _ := os.Hostname()
		response.Hostname = hostname
		response.Secret = viper.GetString("rmq.secret")
	case "agent_stream-init":
		response := q.ExchangeStreamInit{}

		response.StreamName = message["StreamName"]
		response.StreamToken = message["StreamToken"]
		response.JWT = JWT
	}

	// Marshal JSON
	publishMessage, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	// Publish messages to exchange
	publishErr := rmqChan.Publish(
		exchange, // Exchange
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
