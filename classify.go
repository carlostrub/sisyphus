package sisyphus

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
	"github.com/gonum/stat"
	"github.com/retailnext/hllpp"
)

// classificationPrior returns the prior probabilities for good and junk
// classes.
func classificationPrior(db *bolt.DB) (g float64, err error) {

	gTotal, jTotal, err := classificationStatistics(db)
	if err != nil {
		return g, err
	}

	return gTotal / (gTotal + jTotal), err
}

// classificationLikelihoodWordcounts gets wordcounts from database to be used
// in Likelihood calculation
func classificationLikelihoodWordcounts(db *bolt.DB, word string) (gN, jN float64, err error) {

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wordlists"))

		good := b.Bucket([]byte("Good"))
		gWordRaw := good.Get([]byte(word))
		if len(gWordRaw) > 0 {
			var gWordHLL *hllpp.HLLPP
			gWordHLL, err = hllpp.Unmarshal(gWordRaw)
			if err != nil {
				return err
			}
			gN = float64(gWordHLL.Count())
		}
		junk := b.Bucket([]byte("Junk"))
		jWordRaw := junk.Get([]byte(word))
		if len(jWordRaw) > 0 {
			var jWordHLL *hllpp.HLLPP
			jWordHLL, err = hllpp.Unmarshal(jWordRaw)
			if err != nil {
				return err
			}
			jN = float64(jWordHLL.Count())
		}

		return nil
	})

	return gN, jN, err
}

// classificationStatistics gets global statistics from database to
// be used in Likelihood calculation
func classificationStatistics(db *bolt.DB) (gTotal, jTotal float64, err error) {

	err = db.View(func(tx *bolt.Tx) error {
		p := tx.Bucket([]byte("Statistics"))
		gRaw := p.Get([]byte("ProcessedGood"))
		if len(gRaw) > 0 {
			var gHLL *hllpp.HLLPP
			gHLL, err = hllpp.Unmarshal(gRaw)
			if err != nil {
				return err
			}
			gTotal = float64(gHLL.Count())
		}
		jRaw := p.Get([]byte("ProcessedJunk"))
		if len(jRaw) > 0 {
			var jHLL *hllpp.HLLPP
			jHLL, err = hllpp.Unmarshal(jRaw)
			if err != nil {
				return err
			}
			jTotal = float64(jHLL.Count())
		}

		if gTotal == 0 && jTotal == 0 {
			log.Warning("no mails have yet been learned")
			return nil
		}
		if gTotal == 0 {
			log.Warning("no good mails have yet been learned")
			return nil
		}
		if jTotal == 0 {
			log.Warning("no junk mails have yet been learned")
			return nil
		}

		return nil
	})

	return gTotal, jTotal, err
}

// classificationLikelihood returns P(W|C_j) -- the probability of seeing a
// particular word W in a document of this class.
func classificationLikelihood(db *bolt.DB, word string) (g, j float64, err error) {

	gN, jN, err := classificationLikelihoodWordcounts(db, word)
	if err != nil {
		return g, j, err
	}

	gTotal, jTotal, err := classificationStatistics(db)
	if err != nil {
		return g, j, err
	}

	g = gN / gTotal
	j = jN / jTotal

	return g, j, err
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
func (m *Mail) Classify(db *bolt.DB, dir Maildir) (err error) {

	m.New = true

	err = m.Load(dir)
	if err != nil {
		return err
	}

	list, err := m.cleanWordlist()
	if err != nil {
		return err
	}

	junk, prob, err := Junk(db, list)
	if err != nil {
		return err
	}

	m.Junk = junk

	log.WithFields(log.Fields{
		"mail":        m.Key,
		"junk":        m.Junk,
		"probability": prob,
		"dir":         string(dir),
	}).Info("Classified")

	// Move mail around if junk.
	if junk {
		err = os.Rename(filepath.Join(string(dir), "new", m.Key), filepath.Join(string(dir), ".Junk", "cur", m.Key))
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"mail": m.Key,
		}).Info("Moved to Junk folder")
	}

	err = m.Unload(dir)

	return err
}

// Junk returns true if the wordlist is classified as a junk mail using Bayes'
// rule. If required, it also returns the calculated probability of being junk,
// but this is typically not needed.
func Junk(db *bolt.DB, wordlist []string) (junk bool, prob float64, err error) {
	var probabilities []float64

	// initial value should be no junk
	prob = 1.0

	for _, val := range wordlist {
		var p float64
		p, err = classificationWord(db, val)
		if err != nil {
			return false, 0.0, err
		}
		probabilities = append(probabilities, p)
	}

	if len(probabilities) > 0 {
		prob = stat.HarmonicMean(probabilities, nil)
	}
	if prob < 0.5 {
		return true, (1 - prob), err
	}

	return false, (1 - prob), err
}
