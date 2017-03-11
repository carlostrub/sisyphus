package main

import (
	"github.com/boltdb/bolt"
)

// openDB creates and opens a new database and its respective buckets (if required)
func openDB(path string) (db *bolt.DB, err error) {

	// Open the sisyphus.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return db, err
	}

	// Create DB bucket for the map of processed e-mail IDs
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Processed"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return db, err
	}

	// Create DB bucket for word lists
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Wordlists"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return db, err
	}

	// Create DB bucket for Junk inside bucket Wordlists
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wordlists"))
		_, err := b.CreateBucketIfNotExists([]byte("Junk"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return db, err
	}

	// Create DB bucket for Junk inside bucket Wordlists
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wordlists"))
		_, err := b.CreateBucketIfNotExists([]byte("Good"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return db, err
	}

	return db, nil
}
