package sisyphus

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/retailnext/hllpp"
)

// Learn adds the the mail key to the list of words using hyper log log algorithm.
func (m *Mail) Learn(db *bolt.DB) error {

	log.Println("learn mail " + m.Key)

	err := m.Clean()
	if err != nil {
		return err
	}

	list, err := m.Wordlist()
	if err != nil {
		return err
	}

	wordKey := "Good"
	if m.Junk {
		wordKey = "Junk"
	}

	// Learn words
	for _, val := range list {
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Wordlists"))

			bucket := b.Bucket([]byte(wordKey))
			wordRaw := bucket.Get([]byte(val))
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

			err = bucket.Put([]byte(val), word.Marshal())

			return err
		})
		if err != nil {
			return err
		}
	}

	// Update the statistics counter
	err = db.Update(func(tx *bolt.Tx) error {
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
