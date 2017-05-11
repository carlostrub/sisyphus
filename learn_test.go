package sisyphus_test

import (
	"os"

	"github.com/boltdb/bolt"
	. "github.com/carlostrub/sisyphus"
	"github.com/retailnext/hllpp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var m *Mail
var dbs map[Maildir]*bolt.DB
var err error

var _ = Describe("Learn", func() {
	Context("Learn a new mail", func() {

		BeforeEach(func() {
			// check whether there exists a DB file
			_, oserr := os.Stat("test/Maildir/sisyphus.db")
			Ω(os.IsNotExist(oserr)).Should(BeTrue())

			// Load db
			dbs, err = LoadDatabases([]Maildir{"test/Maildir"})
			Ω(err).ShouldNot(HaveOccurred())

			// Load mail
			m = new(Mail)
			m = &Mail{
				Key:     "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
				Subject: nil,
				Body:    nil,
				Junk:    true,
			}

			err = m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			// Cleanup
			CloseDatabases(dbs)

			err = os.Remove("test/Maildir/sisyphus.db")
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("Load databases for a maildir, then learn a mail and check whether the word counts are correct in the db", func() {

			m.Learn(dbs["test/Maildir"])

			var jN, sN, gN int

			err = dbs["test/Maildir"].View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Wordlists"))

				junk := b.Bucket([]byte("Junk"))
				jN = junk.Stats().KeyN

				s := tx.Bucket([]byte("Statistics"))
				sN = s.Stats().KeyN

				return nil
			})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(gN).Should(Equal(0))
			Ω(jN).Should(Equal(27))
			Ω(sN).Should(Equal(1))

		})

		It("Load databases for a maildir, then learn a mail and check whether individual word counts are equal to 1", func() {

			m.Learn(dbs["test/Maildir"])

			var wordCount uint64

			err = dbs["test/Maildir"].View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Wordlists"))

				junk := b.Bucket([]byte("Junk"))

				jWordRaw := junk.Get([]byte("looking"))
				if len(jWordRaw) != 0 {
					jWordHLL, err := hllpp.Unmarshal(jWordRaw)
					if err != nil {
						return err
					}
					wordCount = jWordHLL.Count()
				}

				return nil
			})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(wordCount).Should(Equal(uint64(1)))

		})
	})
})
