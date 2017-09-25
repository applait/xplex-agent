package main

import (
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
)

func rmqDial() *amqp.Channel {
	// Dial connection with RabbitMQ / central virtualhost
	rmqConn, err := amqp.Dial(viper.GetString("rmq.host"))
	if err != nil {
		log.Fatalln("Error establishing connection "+viper.GetString("rmq.host"), err)
	}

	log.Println("Connected to " + viper.GetString("rmq.host"))

	// Open channel with central virtualhost conn
	rmqChan, err := rmqConn.Channel()
	if err != nil {
		log.Fatalln("Error opening channel "+viper.GetString("rmq.host"), err)
	}

	return rmqChan
}
