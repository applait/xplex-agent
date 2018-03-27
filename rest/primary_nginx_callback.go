package rest

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

type StreamEgresses struct {
	Message string `json:"message"`
	Payload []struct {
		ID        int    `json:"id"`
		Service   string `json:"service"`
		StreamKey string `json:"streamKey"`
		RtmpURL   string `json:"rtmpUrl"`
		IsActive  bool   `json:"isActive"`
	} `json:"payload"`
}

// DISABLED. Will be used in future
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
	// DISABLED because rig can't handle connect requests
	// at the moment
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
	//	log.Println(err)
	// }

	// Rig currently sends just 200 OK. Using it to log.
	// That's because nginx-rtmp will not post stream name before on_publish.

	// Proxy back the response status
	// w.WriteHeader(proxyResp.StatusCode)

	w.WriteHeader(http.StatusOK)
}

func rtmpPublish(w http.ResponseWriter, r *http.Request) {
	req := new(rtmpPublishReq)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	err = decoder.Decode(req, r.PostForm)
	if err != nil {
		log.Println(err)
	}

	httpClient := &http.Client{
		Timeout: time.Second * 20,
	}

	rigReq, _ := http.NewRequest("GET", viper.GetString("rig.URL")+"/agents/config/"+req.Name, nil)

	rigRes, err := httpClient.Do(rigReq)
	if err != nil {
		log.Println(err)
	}
	defer rigRes.Body.Close()

	// Store stream specific secondary nginx configuration
	configPath := viper.GetString("nginx.sec.configBasePath") + "/" + req.Name
	_, err = os.Create(configPath)

	streamEgresses := StreamEgresses{}
	err = json.NewDecoder(rigRes.Body).Decode(&streamEgresses)

	if err != nil {
		log.Println(err)
	}

	defer rigRes.Body.Close()

	configStubPath := viper.GetString("nginx.sec.configStubPath")

	var streamEgressURLs []string
	for _, value := range streamEgresses.Payload {
		streamEgressURLs = append(streamEgressURLs, value.RtmpURL)
	}

	// Get available open port for HTTP
	openHTTPPort, err := GetFreePort()
	if err != nil {
		log.Println(err)
	}

	// Get available open port for HTTP
	openRTMPPort, err := GetFreePort()
	if err != nil {
		log.Println(err)
	}

	err = updateSecNginxConf(configPath, configStubPath, streamEgressURLs, openHTTPPort, openRTMPPort)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// PID path
		pidPath := viper.GetString("nginx.sec.pidBasePath") + "/" + req.Name

		// Spin streamer
		err = execworker.StartStreamer("nginx", configPath, pidPath)
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Location", "rtmp://127.0.0.1:"+strconv.Itoa(openRTMPPort)+"/live/"+req.Name)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
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

	// DISABLED. Will be handled when rig enables it
	// proxyResp, err := proxyReq(req)
	// if err != nil {
	//	log.Fatal(err)
	// }

	// PID path
	pidPath := viper.GetString("nginx.sec.pidBasePath") + "/" + req.Name

	// Stop secondary nginx process on stream end
	err = execworker.StopStreamer("nginx", pidPath)
	if err != nil {
		log.Println(err)
	}

	// Proxy back status code to primary nginx
	w.WriteHeader(http.StatusOK)
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
