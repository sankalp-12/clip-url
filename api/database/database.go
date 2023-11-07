package database

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v4"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func CreateClient() (*badger.DB, error) {
	db, err := badger.Open(badger.DefaultOptions("./data"))
	if err != nil {
		fmt.Println("Unable to connect to BadgerDB")
		return nil, err
	}

	return db, err
}

func CreateInfluxClient() influxdb2.Client {
	token := os.Getenv("INFLUXDB_TOKEN")
	url := os.Getenv("INFLUXDB_URL")
	client := influxdb2.NewClient(url, token)
	return client
}
