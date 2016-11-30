package db

import (
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

const DB_PASSCODES = "passcodes"
const DB_CARDS = "cards"

func init() {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DB_PASSCODES))
		_, err = tx.CreateBucketIfNotExists([]byte(DB_CARDS))
		return err
	})
}
