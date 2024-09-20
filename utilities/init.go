package utilities

import (
	"log"
	"time"

	"go.etcd.io/bbolt"
)

var db *bbolt.DB

func handle() *bbolt.DB {
	// Initialize BoltDB
	var err error
	db, err = bbolt.Open("patients.db", 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a bucket for storing patient data
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Patients"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
