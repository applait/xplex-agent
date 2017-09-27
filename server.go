package main

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	_ "net/http/pprof"
	"xplex-agent/es"
	"xplex-agent/rest"
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

	// Dial elasticseatch
	es.Dial(
		viper.GetString("elastic.url"),
		viper.GetString("elastic.index"),
	)

	// Agent HTTP interface for Rig and Nginx
	http.ListenAndServe(":"+viper.GetString("server.port"), rest.Start())
}
