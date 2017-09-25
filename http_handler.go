package main

import (
	"github.com/gorilla/schema"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

type Req struct {
	Call     string `schema:"call"`
	Addr     string `schema:"addr"`
	App      string `schema:"app"`
	FlashVer string `schema:"flashVer"`
	SwfUrl   string `schema:"swfUrl"`
	TcUrl    string `schema:"tcUrl"`
	PageUrl  string `schema:"pageUrl"`
}

// Root handler
func rootHandler(w http.ResponseWriter, r *http.Request, rmqChan *amqp.Channel) {
	req := new(Req)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	decoder := schema.NewDecoder()

	err = decoder.Decode(req, r.PostForm)

	if err != nil {
		log.Println(err)
	}

	// TODO: Pass to RabbitMQ
}
