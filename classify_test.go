package sisyphus_test

import (
	"math"
	"os"

	. "github.com/carlostrub/sisyphus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Classify Mails", func() {
	Context("Classify one word from the mail that was ", func() {
		BeforeEach(func() {
			// check whether there exists a DB file
			_, oserr := os.Stat("test/Maildir/sisyphus.db")
			Ω(os.IsNotExist(oserr)).Should(BeTrue())

			// Load db
			dbs, err = LoadDatabases([]Maildir{
				"test/Maildir",
			})
			Ω(err).ShouldNot(HaveOccurred())

			m = new(Mail)

			// Load junk mail
			m = &Mail{
				Key:  "1488226337.M327833P8269.mail.carlostrub.ch,S=6960,W=7161:2,Sa",
				Junk: true,
			}

			err = m.Learn(dbs["test/Maildir"], "test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			// Load good mail
			m = &Mail{
				Key: "1488230510.M141612P8565.mail.carlostrub.ch,S=5978,W=6119",
			}

			err = m.Learn(dbs["test/Maildir"], "test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			// Cleanup
			CloseDatabases(dbs)

			err = os.Remove("test/Maildir/sisyphus.db")
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("learned before and is junk", func() {

			answer, prob, err := Junk(dbs["test/Maildir"], []string{"london"})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(prob).Should(Equal(1.0))
			Ω(answer).Should(BeTrue())

		})

		It("learned before and is good", func() {

			answer, prob, err := Junk(dbs["test/Maildir"], []string{"localbase"})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(prob).Should(Equal(0.0))
			Ω(answer).Should(BeFalse())

		})

		It("never learned before", func() {

			answer, prob, err := Junk(dbs["test/Maildir"], []string{"abcdefg"})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(math.IsNaN(prob)).Should(BeTrue())
			Ω(answer).Should(BeFalse())

		})

		It("learned both as good and junk, respectively", func() {

			answer, prob, err := Junk(dbs["test/Maildir"], []string{"than"})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(prob).Should(Equal(0.5))
			Ω(answer).Should(BeFalse())

		})
	})

	Context("Do not classify as junk if there is no information", func() {
		BeforeEach(func() {
			// Load empty Maildir2
			err = LoadMaildirs([]Maildir{
				"test/Maildir2",
			})
			Ω(err).ShouldNot(HaveOccurred())

			// Load db
			dbs, err = LoadDatabases([]Maildir{
				"test/Maildir2",
			})
			Ω(err).ShouldNot(HaveOccurred())

		})
		AfterEach(func() {
			// Cleanup
			CloseDatabases(dbs)

			err = os.RemoveAll("test/Maildir2")
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("learned nothing and thus return always good", func() {

			answer, prob, err := Junk(dbs["test/Maildir2"], []string{"Carlo"})

			Ω(err).ShouldNot(HaveOccurred())
			Ω(math.IsNaN(prob)).Should(BeTrue())
			Ω(answer).Should(BeFalse())

		})
	})

	Context("Only classify a random subset of the words in overly long mails", func() {
		BeforeEach(func() {
			// Load empty Maildir2
			err = LoadMaildirs([]Maildir{
				"test/Maildir2",
			})
			Ω(err).ShouldNot(HaveOccurred())

			// Load db
			dbs, err = LoadDatabases([]Maildir{
				"test/Maildir2",
			})
			Ω(err).ShouldNot(HaveOccurred())

		})
		AfterEach(func() {
			// Cleanup
			CloseDatabases(dbs)

			err = os.RemoveAll("test/Maildir2")
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("learned nothing and thus return always good", func() {

			_, _, err := Junk(dbs["test/Maildir2"], []string{
				"Carlo",
				"0",
				"1",
				"2",
				"3",
				"4",
				"5",
				"6",
				"7",
				"8",
				"9",
				"10",
				"11",
				"12",
				"13",
				"14",
				"15",
				"16",
				"17",
				"18",
				"19",
				"20",
				"21",
				"22",
				"23",
				"24",
				"25",
				"26",
				"27",
				"28",
				"29",
				"30",
				"31",
				"32",
				"33",
				"34",
				"35",
				"36",
				"37",
				"38",
				"39",
				"40",
				"41",
				"42",
				"43",
				"44",
				"45",
				"46",
				"47",
				"48",
				"49",
				"50",
				"51",
				"52",
				"53",
				"54",
				"55",
				"56",
				"57",
				"58",
				"59",
			})

			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
