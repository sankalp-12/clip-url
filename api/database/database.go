package database

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

func CreateClient() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions("./data"))
	if err != nil {
		fmt.Println("Couldn't form connection with db")
	}
	return db
}
