package cron

import (
	"net/http"

	"github.com/spf13/viper"
)

func dropStream(streamName string) error {
	dropReq, err := http.NewRequest("GET", viper.GetString("nginx.pri.controlURL")+"drop/client", nil)
	if err != nil {
		return err
	}

	// Add query params
	queryParams := dropReq.URL.Query()
	queryParams.Add("app", viper.GetString("nginx.pri.appName"))
	queryParams.Add("name", streamName)
	dropReq.URL.RawQuery = queryParams.Encode()

	client := &http.Client{}
	_, err = client.Do(dropReq)
	if err != nil {
		return err
	}
	defer dropReq.Body.Close()

	return nil
}

func act(streamName string, action string) error {
	switch action {
	case "drop":
		err := dropStream(streamName)
		if err != nil {
			return err
		}
	default:
		// Do nothing
	}

	return nil
}
