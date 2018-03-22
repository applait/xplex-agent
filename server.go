package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/applait/xplex-agent/rest"

	"github.com/spf13/viper"
)

func init() {
	// Set viper path and read configuration
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

	// Cron polls rig every n minutes with current streams info and execute action based on response
	// DISABLED ATM
	// cron.Start()

	// Rest defines nginx callback handlers and stats endpoints
	rest.Start()

	// Agent HTTP interface for Rig and Nginx
	log.Printf("Agent | Port %s", viper.GetString("server.port"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("server.port"), rest.Start()))
}
