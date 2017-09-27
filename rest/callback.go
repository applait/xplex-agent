package rest

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
)

type rtmpConnectReq struct {
	Call     string `schema:"call"`
	Addr     string `schema:"addr"`
	App      string `schema:"app"`
	FlashVer string `schema:"flashVer"`
	SwfUrl   string `schema:"swfUrl"`
	TcUrl    string `schema:"tcUrl"`
	PageUrl  string `schema:"pageUrl"`
}

type rtmpPlayReq struct {
	Call     string `schema:"call"`
	Addr     string `schema:"addr"`
	ClientId string `schema:"clientid"`
	App      string `schema:"app"`
	FlashVer string `schema:"flashVer"`
	SwfUrl   string `schema:"swfUrl"`
	TcUrl    string `schema:"tcUrl"`
	PageUrl  string `schema:"pageUrl"`
	Name     string `schema:"name"`
}

func rtmpConnect(w http.ResponseWriter, r *http.Request) {
	req := new(rtmpConnectReq)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	decoder := schema.NewDecoder()

	err = decoder.Decode(req, r.PostForm)
	if err != nil {
		log.Println(err)
	}

	// TODO: Ping rig and respond back
}

func rtmpPlay(w http.ResponseWriter, r *http.Request) {
	req := new(rtmpPlayReq)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	decoder := schema.NewDecoder()

	err = decoder.Decode(req, r.PostForm)
	if err != nil {
		log.Println(err)
	}

	// TODO: Ping rig and respond back

}

func rtmpHandler(r *mux.Router) {
	rpost := r.Methods("POST").Subrouter()

	rpost.HandleFunc("/on_connect", rtmpConnect)
	rpost.HandleFunc("/on_play", rtmpPlay)
}

func callbackHandler(r *mux.Router) {
	rtmpHandler(r.PathPrefix("/rtmp").Subrouter())
}
