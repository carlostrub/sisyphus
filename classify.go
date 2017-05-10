package sisyphus

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gonum/stat"
	"github.com/retailnext/hllpp"
)

// classificationPrior returns the prior probabilities for good and junk
// classes.
func classificationPrior(db *bolt.DB) (g float64, err error) {

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wordlists"))

		good := b.Bucket([]byte("Good"))
		gN := float64(good.Stats().KeyN)

		junk := b.Bucket([]byte("Junk"))
		jN := float64(junk.Stats().KeyN)

		// division by zero means there are no learned mails so far
		if (gN + jN) == 0 {
			return errors.New("no mails have been classified so far")
		}

		g = gN / (gN + jN)

		return nil
	})

	return g, err
}

// classificationLikelihood returns P(W|C_j) -- the probability of seeing a
// particular word W in a document of this class.
func classificationLikelihood(db *bolt.DB, word string) (g, j float64, err error) {

	err = db.View(func(tx *bolt.Tx) error {
		var gN, jN uint64

		b := tx.Bucket([]byte("Wordlists"))

		good := b.Bucket([]byte("Good"))
		gWordRaw := good.Get([]byte(word))
		if len(gWordRaw) != 0 {
			gWordHLL, err := hllpp.Unmarshal(gWordRaw)
			if err != nil {
				return err
			}
			gN = gWordHLL.Count()
		}
		junk := b.Bucket([]byte("Junk"))
		jWordRaw := junk.Get([]byte(word))
		if len(jWordRaw) != 0 {
			jWordHLL, err := hllpp.Unmarshal(jWordRaw)
			if err != nil {
				return err
			}
			jN = jWordHLL.Count()
		}

		p := tx.Bucket([]byte("Statistics"))
		gHLL, err := hllpp.Unmarshal(p.Get([]byte("ProcessedGood")))
		if err != nil {
			return err
		}
		jHLL, err := hllpp.Unmarshal(p.Get([]byte("ProcessedJunk")))
		if err != nil {
			return err
		}

		gTotal := gHLL.Count()
		if gTotal == 0 {
			return errors.New("no good mails have been classified so far")
		}
		jTotal := jHLL.Count()
		if jTotal == 0 {
			return errors.New("no junk mails have been classified so far")
		}

		g = float64(gN) / float64(gTotal)
		j = float64(jN) / float64(jTotal)

		return nil
	})

	return g, j, nil
}

// classificationWord produces the conditional probability of a word belonging
// to good or junk using the classic Bayes' rule.
func classificationWord(db *bolt.DB, word string) (g float64, err error) {

	priorG, err := classificationPrior(db)
	if err != nil {
		return g, err
	}

	likelihoodG, likelihoodJ, err := classificationLikelihood(db, word)
	if err != nil {
		return g, err
	}

	g = (likelihoodG * priorG) / (likelihoodG*priorG + likelihoodJ*(1-priorG))

	return g, nil
}

// Classify analyses a new mail (a mail that arrived in the "new" directory),
// decides whether it is junk and -- if so -- moves it to the Junk folder. If
// it is not junk, the mail is untouched so it can be handled by the mail
// client.
func (m *Mail) Classify(db *bolt.DB) error {

	err := m.Clean()
	if err != nil {
		return err
	}

	list := m.Wordlist()
	junk, err := Junk(db, list)
	if err != nil {
		return err
	}

	log.Print("Classified " + m.Key + " as Junk=" + strconv.FormatBool(m.Junk))

	// Move mail around if junk.
	if junk {
		m.Junk = junk
		err := os.Rename("./new/"+m.Key, "./.Junk/cur/"+m.Key)
		if err != nil {
			return err
		}
		log.Print("Moved " + m.Key + " from new to Junk folder")
	}

	return nil
}

// Junk returns true if the wordlist is classified as a junk mail using Bayes'
// rule.
func Junk(db *bolt.DB, wordlist []string) (bool, error) {
	var probabilities []float64

	for _, val := range wordlist {
		p, err := classificationWord(db, val)
		if err != nil {
			return false, err
		}
		probabilities = append(probabilities, p)
	}

	if stat.HarmonicMean(probabilities, nil) < 0.5 {
		return true, nil
	}

	return false, nil
}
