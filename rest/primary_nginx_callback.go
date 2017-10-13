package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/applait/xplex-agent/execworker"
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

type rtmpPublishReq struct {
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

func proxyReq(req interface{}) (*http.Response, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
	}

	proxyReq, err := http.NewRequest("POST", viper.GetString("rig.URL"), bytes.NewBuffer(payload))
	proxyReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	proxyResp, err := client.Do(proxyReq)
	if err != nil {
		return proxyResp, err
	}
	defer proxyReq.Body.Close()

	return proxyResp, nil
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

	proxyResp, err := proxyReq(req)
	if err != nil {
		log.Fatal(err)
	}

	// Proxy back status code to primary nginx
	w.WriteHeader(proxyResp.StatusCode)
}

func rtmpPublish(w http.ResponseWriter, r *http.Request) {
	req := new(rtmpPublishReq)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	decoder := schema.NewDecoder()

	err = decoder.Decode(req, r.PostForm)
	if err != nil {
		log.Println(err)
	}

	proxyResp, err := proxyReq(req)
	if err != nil {
		log.Fatal(err)
	}

	// Store stream specific secondary nginx configuration
	out, err := os.Create(viper.GetString("nginx.streamConfigBasePath") + req.Name)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	io.Copy(out, proxyResp.Body)

	// Spin secondary nginx process
	execworker.SpinNginx(req.Name)

	// Proxy back status code to primary nginx
	w.WriteHeader(proxyResp.StatusCode)
}

func rtmpHandler(r *mux.Router) {
	rpost := r.Methods("POST").Subrouter()

	rpost.HandleFunc("/on_connect", rtmpConnect)
	rpost.HandleFunc("/on_publish", rtmpPublish)
}

func callbackHandler(r *mux.Router) {
	rtmpHandler(r.PathPrefix("/rtmp").Subrouter())
}
