package sisyphus_test

import (
	"math"
	"os"

	. "github.com/carlostrub/sisyphus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Classify", func() {
	Context("Classify a new mail", func() {
		BeforeEach(func() {
			// check whether there exists a DB file
			_, oserr := os.Stat("test/Maildir/sisyphus.db")
			Ω(os.IsNotExist(oserr)).Should(BeTrue())

			// Load db
			dbs, err = LoadDatabases([]Maildir{"test/Maildir"})
			Ω(err).ShouldNot(HaveOccurred())

			m = new(Mail)

			// Load junk mail
			m = &Mail{
				Key:  "1488226337.M327833P8269.mail.carlostrub.ch,S=6960,W=7161:2,Sa",
				Junk: true,
			}

			err = m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Learn(dbs["test/Maildir"])
			Ω(err).ShouldNot(HaveOccurred())

			// Load good mail
			m = &Mail{
				Key: "1488230510.M141612P8565.mail.carlostrub.ch,S=5978,W=6119",
			}

			err = m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Learn(dbs["test/Maildir"])
			Ω(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			// Cleanup
			CloseDatabases(dbs)

			err = os.Remove("test/Maildir/sisyphus.db")
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("Classify one word from the mail that was learned before", func() {

			answer, prob, err := Junk(dbs["test/Maildir"], []string{"london"})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(prob).Should(Equal(1.0))
			Ω(answer).Should(BeTrue())

		})

		It("Classify one word from the mail that was learned before", func() {

			answer, prob, err := Junk(dbs["test/Maildir"], []string{"localbase"})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(prob).Should(Equal(0.0))
			Ω(answer).Should(BeFalse())

		})

		It("Classify one word from the mail that was never learned", func() {

			answer, prob, err := Junk(dbs["test/Maildir"], []string{"abcdefg"})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(math.IsNaN(prob)).Should(BeTrue())
			Ω(answer).Should(BeFalse())

		})

		It("Classify one word from the mail that was learned in good and junk", func() {

			answer, prob, err := Junk(dbs["test/Maildir"], []string{"than"})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(prob).Should(Equal(0.7795275590551181))
			Ω(answer).Should(BeTrue())

		})
	})
})
