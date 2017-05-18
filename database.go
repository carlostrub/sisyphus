package sisyphus

import (
	log "github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
)

// openDB creates and opens a new database and its respective buckets (if required)
func openDB(m Maildir) (db *bolt.DB, err error) {

	log.Println("loading database for " + string(m))
	// Open the sisyphus.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err = bolt.Open(string(m)+"/sisyphus.db", 0600, nil)
	if err != nil {
		return db, err
	}

	// Create DB bucket for the map of processed e-mail IDs
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Statistics"))
		return err
	})
	if err != nil {
		return db, err
	}

	// Create DB bucket for word lists
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Wordlists"))
		return err
	})
	if err != nil {
		return db, err
	}

	// Create DB bucket for Junk inside bucket Wordlists
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wordlists"))
		_, err := b.CreateBucketIfNotExists([]byte("Junk"))
		return err
	})
	if err != nil {
		return db, err
	}

	// Create DB bucket for Good inside bucket Wordlists
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wordlists"))
		_, err := b.CreateBucketIfNotExists([]byte("Good"))
		return err
	})

	return db, err
}

// LoadDatabases loads all databases from a given slice of Maildirs
func LoadDatabases(d []Maildir) (databases map[Maildir]*bolt.DB, err error) {
	databases = make(map[Maildir]*bolt.DB)
	for _, val := range d {
		databases[val], err = openDB(val)
		if err != nil {
			return databases, err
		}
	}

	log.Println("all databases loaded")

	return databases, nil
}

// CloseDatabases closes all databases from a given slice of Maildirs
func CloseDatabases(databases map[Maildir]*bolt.DB) {
	for key, val := range databases {
		err := val.Close()
		if err != nil {
			log.Println(err)
		}
		log.Println("database " + string(key) + "/sisyphus.db closed")
	}

	return
}
