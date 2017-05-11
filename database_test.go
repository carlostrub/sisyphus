package sisyphus_test

import (
	"os"

	"github.com/boltdb/bolt"
	. "github.com/carlostrub/sisyphus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {

	Context("Bolt Database", func() {
		AfterEach(func() {
			err = os.Remove("test/Maildir/sisyphus.db")
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("Load databases for each maildir", func() {
			dbs, err := LoadDatabases([]Maildir{"test/Maildir"})
			Ω(err).ShouldNot(HaveOccurred())

			dbTest := dbs["test/Maildir"]
			var gN = 4
			var jN = 4
			var sN = 4

			err = dbTest.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Wordlists"))

				good := b.Bucket([]byte("Good"))
				gN = good.Stats().KeyN

				junk := b.Bucket([]byte("Junk"))
				jN = junk.Stats().KeyN

				s := tx.Bucket([]byte("Statistics"))
				sN = s.Stats().KeyN

				return nil
			})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(gN).Should(Equal(0))
			Ω(jN).Should(Equal(0))
			Ω(sN).Should(Equal(0))

			CloseDatabases(dbs)
		})

		It("Closes an open database", func() {
			dbs, err := LoadDatabases([]Maildir{"test/Maildir"})
			Ω(err).ShouldNot(HaveOccurred())

			CloseDatabases(dbs)

			dbTest := dbs["test/Maildir"]
			var n = 4

			err = dbTest.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Statistics"))

				n = b.Stats().KeyN

				return nil
			})
			Ω(err).Should(HaveOccurred())
		})
	})
})
