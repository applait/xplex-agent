package main

import (
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var JWT = ""

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	// Handle errors reading the config file
	if err != nil {
		log.Fatalln("Fatal error config file", err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Rabbitmq: Dial
	rmqChan := rmqDial()

	// Rabbitmq: Declare exchange - agent-register
	err := rmqChan.ExchangeDeclare(
		"agent_register", // name of the exchange
		"direct",         // type
		true,             // durable
		false,            // delete when complete
		false,            // internal
		false,            // noWait
		nil,              // arguments
	)

	if err != nil {
		log.Fatalln("Error declaring exchange: agent_register", err)
	}

	// Rabbitmq: Declare exchange - agent-register
	err = rmqChan.ExchangeDeclare(
		"agent_stream-init", // name of the exchange
		"direct",            // type
		true,                // durable
		false,               // delete when complete
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	)

	if err != nil {
		log.Fatalln("Error declaring exchange: agent_stream-init", err)
	}

	// HTTP: Router
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rootHandler(w, r, rmqChan)
	}).Methods("POST")

	// HTTP: Listener
	http.ListenAndServe(":"+viper.GetString("server.port"), router)
}
