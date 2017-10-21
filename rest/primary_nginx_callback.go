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
	// req := new(rtmpConnectReq)
	// err := r.ParseForm()
	// if err != nil {
	//	log.Println(err)
	// }

	// decoder := schema.NewDecoder()

	// err = decoder.Decode(req, r.PostForm)
	// if err != nil {
	//	log.Println(err)
	// }

	// proxyResp, err := proxyReq(req)
	// if err != nil {
	//	log.Fatal(err)
	// }

	// Returning 200 for time being
	w.WriteHeader(http.StatusOK)
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
	configPath := viper.GetString("nginx.sec.configBasePath") + req.Name
	out, err := os.Create(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	io.Copy(out, proxyResp.Body)

	// PID path
	pidPath := viper.GetString("nginx.sec.pidBasePath") + req.Name

	// Spin secondary nginx process
	startNginxErr := execworker.StartNginx(configPath, pidPath)
	if startNginxErr != nil {
		log.Println(err)
	}

	// Proxy back status code to primary nginx
	w.WriteHeader(proxyResp.StatusCode)
}

func rtmpPublishDone(w http.ResponseWriter, r *http.Request) {
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

	// PID path
	pidPath := viper.GetString("nginx.sec.pidBasePath") + req.Name

	// Stop secondary nginx process on stream end
	stopNginxErr := execworker.StopNginx(pidPath)
	if stopNginxErr != nil {
		log.Println(err)
	}

	// Proxy back status code to primary nginx
	w.WriteHeader(proxyResp.StatusCode)
}

func rtmpHandler(r *mux.Router) {
	rpost := r.Methods("POST").Subrouter()

	// Connect start handler. Used for authentication purposes.
	rpost.HandleFunc("/on_connect", rtmpConnect)

	// Publish start/end callback handlers
	rpost.HandleFunc("/on_publish", rtmpPublish)
	rpost.HandleFunc("/on_publish_done", rtmpPublishDone)
}

func callbackHandler(r *mux.Router) {
	rtmpHandler(r.PathPrefix("/rtmp").Subrouter())
}
