package database

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

func CreateClient() (*badger.DB, error) {
	db, err := badger.Open(badger.DefaultOptions("./data"))
	if err != nil {
		fmt.Println("Unable to connect to BadgerDB")
		return nil, err
	}

	return db, err
}
