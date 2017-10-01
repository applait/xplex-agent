package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"xplex-agent/cron"
	"xplex-agent/rest"

	"github.com/spf13/viper"
)

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

	// Poll Rig on stream status
	cron.Start()

	// HTTP route handler
	rest.Start()

	// Agent HTTP interface for Rig and Nginx
	http.ListenAndServe(":"+viper.GetString("server.port"), rest.Start())
}
