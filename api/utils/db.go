package utils

import (
	"fmt"

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

func WriteToDB(db *badger.DB, oldURL string, newURL string) (bool, string) {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	key := []byte(fmt.Sprintf("w:%s", oldURL))
	value := []byte(newURL)
	err := txn.Set(key, value)
	if err != nil {
		return false, "Unable to write to the DB"
	}

	key = []byte(fmt.Sprintf("r:%s", newURL))
	value = []byte(oldURL)
	err = txn.Set(key, value)
	if err != nil {
		return false, "Unable to write to the DB"
	}

	if err = txn.Commit(); err != nil {
		return false, "Unable to commit transaction to the DB"
	}

	return true, "Write successful"
}
