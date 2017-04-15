/*
Part of this code is borrowed from github.com/jbrukh/bayesian published under a BSD3CLAUSE License
*/

package sisyphus

import (
	"math"
	"strconv"

	"github.com/boltdb/bolt"
)

// classificationPriors returns the prior probabilities for good and junk
// classes.
func classificationPriors(db *bolt.DB) (g, j float64) {

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wordlists"))
		good := b.Bucket([]byte("Good"))
		gN := float64(good.Stats().KeyN)
		junk := b.Bucket([]byte("Junk"))
		jN := float64(junk.Stats().KeyN)

		g = gN / (gN + jN)
		j = jN / (gN + jN)

		return nil
	})

	return
}

// classificationWordProb returns P(W|C_j) -- the probability of seeing
// a particular word W in a document of this class.
func classificationWordProb(db *bolt.DB, word string) (g, j float64) {

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Wordlists"))
		good := b.Bucket([]byte("Good"))
		gNString := string(good.Get([]byte(word)))
		gN, _ := strconv.ParseFloat(gNString, 64)
		junk := b.Bucket([]byte("Junk"))
		jNString := string(junk.Get([]byte(word)))
		jN, _ := strconv.ParseFloat(jNString, 64)

		p := tx.Bucket([]byte("Processed"))
		counters := p.Bucket([]byte("Counters"))
		jString := string(counters.Get([]byte("Junk")))
		j, _ := strconv.ParseFloat(jString, 64)
		mails := p.Bucket([]byte("Mails"))
		pN := mails.Stats().KeyN

		g = gN / (float64(pN) - j)
		j = jN / j

		return nil
	})

	return g, j
}

// LogScores produces "log-likelihood"-like scores that can
// be used to classify documents into classes.
//
// The value of the score is proportional to the likelihood,
// as determined by the classifier, that the given document
// belongs to the given class. This is true even when scores
// returned are negative, which they will be (since we are
// taking logs of probabilities).
//
// The index j of the score corresponds to the class given
// by c.Classes[j].
//
// Additionally returned are "inx" and "strict" values. The
// inx corresponds to the maximum score in the array. If more
// than one of the scores holds the maximum values, then
// strict is false.
//
// Unlike c.Probabilities(), this function is not prone to
// floating point underflow and is relatively safe to use.
func LogScores(db *bolt.DB, wordlist []string) (scoreG, scoreJ float64, junk bool) {

	priorG, priorJ := classificationPriors(db)

	// calculate the scores
	scoreG = math.Log(priorG)
	scoreJ = math.Log(priorJ)
	for _, word := range wordlist {
		gP, jP := classificationWordProb(db, word)
		scoreG += math.Log(gP)
		scoreJ += math.Log(jP)
	}

	if scoreJ == math.Max(scoreG, scoreJ) {
		junk = true
	}

	return scoreG, scoreJ, junk
}
