package sisyphus

import (
	log "github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
	"github.com/retailnext/hllpp"
)

// learnWordlist adds the mail key to the respective word's list
func (m *Mail) learnWordlist(w string, db *bolt.DB) error {
	wordKey := "Good"
	if m.Junk {
		wordKey = "Junk"
	}

	err := db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte("Wordlists"))

		bucket := b.Bucket([]byte(wordKey))
		wordRaw := bucket.Get([]byte(w))
		var word *hllpp.HLLPP
		if len(wordRaw) == 0 {
			word = hllpp.New()
		} else {
			word, err = hllpp.Unmarshal(wordRaw)
			if err != nil {
				return err
			}
		}

		word.Add([]byte(m.Key))

		err = bucket.Put([]byte(w), word.Marshal())

		return err
	})

	return err
}

// learnStatistics adds the mail key to the respective word's list
func (m *Mail) learnStatistics(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) (err error) {
		p := tx.Bucket([]byte("Statistics"))

		key := "ProcessedGood"
		if m.Junk {
			key = "ProcessedJunk"
		}

		keyRaw := p.Get([]byte(key))
		var counter *hllpp.HLLPP
		if len(keyRaw) == 0 {
			counter = hllpp.New()
		} else {
			counter, err = hllpp.Unmarshal(keyRaw)
			if err != nil {
				return err
			}
		}

		counter.Add([]byte(m.Key))

		err = p.Put([]byte(key), counter.Marshal())

		return err
	})

	return err
}

// Learn adds the the mail key to the list of words using hyper log log algorithm.
func (m *Mail) Learn(db *bolt.DB) error {

	log.WithFields(log.Fields{
		"mail": m.Key,
	}).Info("Learn mail")

	list, err := m.cleanWordlist()
	if err != nil {
		return err
	}

	// Learn words
	for _, val := range list {
		err := m.learnWordlist(val, db)
		if err != nil {
			return err
		}
	}

	// Update the statistics counter
	err = m.learnStatistics(db)

	return err

}
