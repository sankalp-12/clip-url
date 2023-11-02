package utils

import "github.com/dgraph-io/badger/v4"

func CheckCollisions(db *badger.DB, prefix []byte, url string) (bool, string) {
	txn := db.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		key := string(item.Key())
		val, _ := item.ValueCopy(nil)

		if key == url {
			return true, string(val)
		}
	}
	defer it.Close()

	return false, ""
}
