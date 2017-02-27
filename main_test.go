package main_test

import (
	. "github.com/carlostrub/sisyphus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {

	Context("Index Maildir", func() {
		It("Create a slice of mail keys", func() {
			result, err := Index("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(Equal(
				[]Mail{

					{
						Key:     "1488230510.M141612P8565.mail.carlostrub.ch,S=5978,W=6119",
						Subject: nil,
						Body:    nil,
						Junk:    false,
					},
					{
						Key:     "1488181583.M633084P4781.mail.carlostrub.ch,S=708375,W=720014",
						Subject: nil,
						Body:    nil,
						Junk:    true,
					},
					{
						Key:     "1488226337.M327824P8269.mail.carlostrub.ch,S=8044,W=8167",
						Subject: nil,
						Body:    nil,
						Junk:    true,
					},
					{
						Key:     "1488226337.M327825P8269.mail.carlostrub.ch,S=802286,W=812785",
						Subject: nil,
						Body:    nil,
						Junk:    true,
					},
					{
						Key:     "1488228352.M339670P8269.mail.carlostrub.ch,S=12659,W=12782",
						Subject: nil,
						Body:    nil,
						Junk:    true,
					},
					{
						Key:     "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
						Subject: nil,
						Body:    nil,
						Junk:    true,
					},
					{
						Key:     "1488226337.M327833P8269.mail.carlostrub.ch,S=6960,W=7161",
						Subject: nil,
						Body:    nil,
						Junk:    true,
					},
				}))
		})
	})
})
