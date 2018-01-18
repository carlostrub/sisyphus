package sisyphus

import (
	"github.com/boltdb/bolt"
	"github.com/retailnext/hllpp"
)

// Info produces statistics
func Info(db *bolt.DB) (gTotal, jTotal, gWords, jWords uint64) {

	_ = db.View(func(tx *bolt.Tx) error {
		p := tx.Bucket([]byte("Statistics"))
		gRaw := p.Get([]byte("ProcessedGood"))
		if len(gRaw) > 0 {
			var gHLL *hllpp.HLLPP
			gHLL, _ = hllpp.Unmarshal(gRaw)
			gTotal = gHLL.Count()
		}
		jRaw := p.Get([]byte("ProcessedJunk"))
		if len(jRaw) > 0 {
			var jHLL *hllpp.HLLPP
			jHLL, _ = hllpp.Unmarshal(jRaw)
			jTotal = jHLL.Count()
		}

		return nil
	})

	_ = db.View(func(tx *bolt.Tx) error {
		p := tx.Bucket([]byte("Wordlists"))
		pj := p.Bucket([]byte("Junk"))

		stats := pj.Stats()
		jWords = uint64(stats.KeyN)

		return nil
	})

	_ = db.View(func(tx *bolt.Tx) error {
		p := tx.Bucket([]byte("Wordlists"))
		pg := p.Bucket([]byte("Good"))

		stats := pg.Stats()
		gWords = uint64(stats.KeyN)

		return nil
	})

	return gTotal, jTotal, gWords, jWords
}
