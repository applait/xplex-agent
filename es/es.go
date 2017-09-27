package es

import (
	"context"
	"github.com/spf13/viper"
	elastic "gopkg.in/olivere/elastic.v5"
)

func Dial(url string, index string) {
	// Create a client and connect to ElasticSearch instance
	client, err := elastic.NewClient(elastic.SetURL(viper.GetString(url)))
	if err != nil {
		panic(err)
	}

	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		// TODO: Create index
	}

}
