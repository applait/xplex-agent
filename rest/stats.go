package rest

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type NginxStats struct {
	Rtmp struct {
		NginxVersion     string `xml:"nginx_version"`
		NginxRtmpVersion string `xml:"nginx_rtmp_version"`
		Compiler         string `xml:"compiler"`
		Built            string `xml:"built"`
		Pid              string `xml:"pid"`
		Uptime           string `xml:"uptime"`
		Naccepted        string `xml:"naccepted"`
		BwIn             string `xml:"bw_in"`
		BytesIn          string `xml:"bytes_in"`
		BwOut            string `xml:"bw_out"`
		BytesOut         string `xml:"bytes_out"`
		Server           struct {
			Application struct {
				Name string `xml:"name"`
				Live struct {
					Stream []struct {
						Name     string `xml:"name"`
						Time     string `xml:"time"`
						BwIn     string `xml:"bw_in"`
						BytesIn  string `xml:"bytes_in"`
						BwOut    string `xml:"bw_out"`
						BytesOut string `xml:"bytes_out"`
						BwAudio  string `xml:"bw_audio"`
						BwVideo  string `xml:"bw_video"`
						Client   []struct {
							ID         string `xml:"id"`
							Address    string `xml:"address"`
							Time       string `xml:"time"`
							Flashver   string `xml:"flashver"`
							Dropped    string `xml:"dropped"`
							Avsync     string `xml:"avsync"`
							Timestamp  string `xml:"timestamp"`
							Active     string `xml:"active"`
							Swfurl     string `xml:"swfurl,omitempty"`
							Publishing string `xml:"publishing,omitempty"`
						} `xml:"client"`
						Meta struct {
							Video struct {
								Width     string `xml:"width"`
								Height    string `xml:"height"`
								FrameRate string `xml:"frame_rate"`
								Codec     string `xml:"codec"`
								Profile   string `xml:"profile"`
								Compat    string `xml:"compat"`
								Level     string `xml:"level"`
							} `xml:"video"`
							Audio struct {
								Codec      string `xml:"codec"`
								Profile    string `xml:"profile"`
								Channels   string `xml:"channels"`
								SampleRate string `xml:"sample_rate"`
							} `xml:"audio"`
						} `xml:"meta"`
						Nclients   string `xml:"nclients"`
						Publishing string `xml:"publishing"`
						Active     string `xml:"active"`
					} `xml:"stream"`
					Nclients string `xml:"nclients"`
				} `xml:"live"`
			} `xml:"application"`
		} `xml:"server"`
	} `xml:"rtmp"`
	Host string `json:"host"`
}

// GetNginxRTMPStats returns a JSON formatted RTMP stats by default
func GetNginxRTMPStats() ([]byte, error) {
	statsRes, err := http.Get(viper.GetString("nginx.pri.statsURL"))
	if err != nil {
		return nil, err
	}
	defer statsRes.Body.Close()

	statsBody, err := ioutil.ReadAll(statsRes.Body)
	if err != nil {
		return nil, err
	}

	nginxStats := new(NginxStats)
	err = xml.Unmarshal(statsBody, &nginxStats)
	if err != nil {
		return nil, err
	}

	nginxStats.Host, err = os.Hostname()

	nginxStatsJSON, err := json.Marshal(nginxStats)
	if err != nil {
		return nil, err
	}

	return nginxStatsJSON, nil
}

// Queries HTTP interface provided by Nginx for stats and returns the response
func streamHandler(w http.ResponseWriter, r *http.Request) {
	nginxRTMPStats, err := GetNginxRTMPStats()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(nginxRTMPStats)
	}
}

func statsHandler(r *mux.Router) {
	rget := r.Methods("GET").Subrouter()

	rget.HandleFunc("/stream/rtmp", streamHandler)
}
