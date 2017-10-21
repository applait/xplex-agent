package cron

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/applait/xplex-agent/rest"
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"
)

type Action struct {
	Streams []struct {
		Name   string `json:"name"`
		Action string `json:"action"`
	} `json:"streams"`
}

func pollRig() {
	// Get nginx RTMP stats
	nginxRTMPStats, err := rest.GetNginxRTMPStats()
	if err != nil {
		log.Println(err)
	}

	// Ping rig with the info and wait for action
	actionReq, err := http.NewRequest("POST", viper.GetString("rig.URL")+"action", bytes.NewBuffer(nginxRTMPStats))
	actionReq.Header.Set("Content-Type", "application/json")

	// Make request
	client := &http.Client{}
	actionResp, err := client.Do(actionReq)
	if err != nil {
		log.Println(err)
	}

	// Parse body
	actionBody, err := ioutil.ReadAll(actionResp.Body)
	if err != nil {
		log.Println(err)
	}

	defer actionReq.Body.Close()

	action := new(Action)
	err = json.Unmarshal(actionBody, &action)
	if err != nil {
		log.Println(err)
	}

	// Loop through response and act
	for _, stream := range action.Streams {
		err = act(stream.Name, stream.Action)
		if err != nil {
			log.Println(err)
		}
	}

}

// Start starts the CRON which polls rig and take necessary actions on streams
func Start() {
	gocron.Every(5).Minutes().Do(pollRig)

	gocron.Start()
}
