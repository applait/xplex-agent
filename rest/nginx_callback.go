package rest

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/spf13/viper"
)

type rtmpConnectReq struct {
	Call     string `schema:"call"`
	Addr     string `schema:"addr"`
	App      string `schema:"app"`
	FlashVer string `schema:"flashVer"`
	SwfURL   string `schema:"swfUrl"`
	TcURL    string `schema:"tcUrl"`
	PageURL  string `schema:"pageUrl"`
}

type rtmpPlayReq struct {
	Call     string `schema:"call"`
	Addr     string `schema:"addr"`
	ClientID string `schema:"clientid"`
	App      string `schema:"app"`
	FlashVer string `schema:"flashVer"`
	SwfURL   string `schema:"swfUrl"`
	TcURL    string `schema:"tcUrl"`
	PageURL  string `schema:"pageUrl"`
	Name     string `schema:"name"`
}

func proxyReq(req interface{}, w http.ResponseWriter) {
	payload, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
	}

	proxyReq, err := http.NewRequest("POST", viper.GetString("rig.URL"), bytes.NewBuffer(payload))
	proxyReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	proxyResp, err := client.Do(proxyReq)
	if err != nil {
		panic(err)
	}
	defer proxyReq.Body.Close()

	w.WriteHeader(proxyResp.StatusCode)
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

	proxyReq(req, w)
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

	proxyReq(req, w)
}

func rtmpHandler(r *mux.Router) {
	rpost := r.Methods("POST").Subrouter()

	rpost.HandleFunc("/on_connect", rtmpConnect)
	rpost.HandleFunc("/on_play", rtmpPlay)
}

func callbackHandler(r *mux.Router) {
	rtmpHandler(r.PathPrefix("/rtmp").Subrouter())
}
