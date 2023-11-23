package utils

import (
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v4"
)

func CheckCollisions(db *badger.DB, prefix []byte, url string) (bool, string) {
	txn := db.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		key := string(item.Key())[2:]
		val, _ := item.ValueCopy(nil)

		if key == url {
			return true, string(val)
		}
	}

	return false, ""
}

func ReturnAllData(db *badger.DB) (bool, map[string]string) {

	allValues := make(map[string]string)
	err := db.View(func(txn *badger.Txn) error {
		iterator := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iterator.Close()

		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			item := iterator.Item()

			// Get the key
			key := item.Key()

			// Get the value
			var value []byte
			err := item.Value(func(val []byte) error {
				value = append([]byte{}, val...)
				return nil
			})
			if err != nil {
				return err
			}
			if string(key)[0] == 'r' {
				// Process the key-value pair
				allValues[strings.Split(string(key), ":")[1]] = string(value)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Some issue here")
		return false, map[string]string{}
	}
	return true, allValues
}

func WriteToDB(db *badger.DB, oldURL string, newURL string) (bool, string) {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	key := []byte(fmt.Sprintf("w:%s", oldURL))
	value := []byte(newURL)
	err := txn.Set(key, value)
	if err != nil {
		return false, "Unable to write to BadgerDB"
	}

	key = []byte(fmt.Sprintf("r:%s", newURL))
	value = []byte(oldURL)
	err = txn.Set(key, value)
	if err != nil {
		return false, "Unable to write to BadgerDB"
	}

	if err = txn.Commit(); err != nil {
		return false, "Unable to commit transaction to BadgerDB"
	}

	return true, "Successfully written to BadgerDB"
}
