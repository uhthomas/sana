package db

import (
	"log"

	elastic "gopkg.in/olivere/elastic.v3"
	uppdb "upper.io/db"
	"upper.io/db/mongo"
)

var (
	Session       uppdb.Database
	ElasticClient *elastic.Client
)

func init() {
	var err error
	Session, err = uppdb.Open(mongo.Adapter, mongo.ConnectionURL{
		Address:  uppdb.Host("127.0.0.1"),
		Database: "sana",
	})
	if err != nil {
		log.Fatal(err)
	}
	ElasticClient, err = elastic.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	exists, err := ElasticClient.IndexExists("sana").Do()
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		_, err := ElasticClient.CreateIndex("sana").Do()
		if err != nil {
			log.Fatal(err)
		}
	}
}
