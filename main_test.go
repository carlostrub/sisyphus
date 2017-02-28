package main_test

import (
	. "github.com/carlostrub/sisyphus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {

	Context("Maildir", func() {
		It("Create a slice of mail keys", func() {
			result, err := Index("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(result).Should(Equal(
				[]*Mail{

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

	Context("Mail", func() {
		It("Load mail content into struct", func() {
			m := Mail{
				Key:     "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
				Subject: nil,
				Body:    nil,
				Junk:    true,
			}

			err := m.Load("test/Maildir" + "/.Junk")
			Ω(err).ShouldNot(HaveOccurred())

			subject := "hello"
			body := "This is a multi-part message in MIME format.  ------=_NextPart_000_0032_01D2912F.05324BC6 Content-Type: text/plain; \tcharset=\"cp-850\" Content-Transfer-Encoding: quoted-printable  Dear cs,  We are looking for employees working remotely.  My name is Kari, I am the personnel manager of a large International company. Most of the work you can do from home, that is, at a distance.  Salary is $2000-$5300.  If you are interested in this offer, please visit  Our Site  Best regards! ------=_NextPart_000_0032_01D2912F.05324BC6 Content-Type: text/html; \tcharset=\"cp-850\" Content-Transfer-Encoding: quoted-printable  <html xmlns:v=\"urn:schemas-microsoft-com:vml\" xmlns:o=\"urn:schemas-microsoft-com:office:office\" xmlns:w=\"urn:schemas-microsoft-com:office:word\" xmlns:m=\"http://schemas.microsoft.com/office/2004/12/omml\" xmlns=\"http://www.w3.org/TR/REC-html40\"><head><META HTTP-EQUIV=\"Content-Type\" CONTENT=\"text/html; charset=us-ascii\"><meta name=Generator content=\"Microsoft Word 14 (filtered medium)\"><style><!-- /* Font Definitions */ @font-face \t{font-family:Calibri; \tpanose-1:2 15 5 2 2 2 4 3 2 4;} /* Style Definitions */ p.MsoNormal, li.MsoNormal, div.MsoNormal \t{margin:0in; \tmargin-bottom:.0001pt; \tfont-size:11.0pt; \tfont-family:\"Calibri\",\"sans-serif\";} a:link, span.MsoHyperlink \t{mso-style-priority:99; \tcolor:blue; \ttext-decoration:underline;} a:visited, span.MsoHyperlinkFollowed \t{mso-style-priority:99; \tcolor:purple; \ttext-decoration:underline;} span.EmailStyle17 \t{mso-style-type:personal-compose; \tfont-family:\"Calibri\",\"sans-serif\"; \tcolor:windowtext;} .MsoChpDefault \t{mso-style-type:export-only; \tfont-family:\"Calibri\",\"sans-serif\";} @page WordSection1 \t{size:8.5in 11.0in; \tmargin:1.0in 1.0in 1.0in 1.0in;} div.WordSection1 \t{page:WordSection1;} --></style><!--[if gte mso 9]><xml> <o:shapedefaults v:ext=\"edit\" spidmax=\"1026\" /> </xml><![endif]--><!--[if gte mso 9]><xml> <o:shapelayout v:ext=\"edit\"> <o:idmap v:ext=\"edit\" data=\"1\" /> </o:shapelayout></xml><![endif]--></head><body lang=EN-US link=blue vlink=purple><div class=WordSection1><p class=MsoNormal>Dear cs,<br> <br> We are looking for employees working remotely.<br> <br> My name is Kari, I am the personnel manager of a large International company.<br> Most of the work you can do from home, that is, at a distance.<br> <b>Salary is $2000-$5300.</b><br> <br> If you are interested in this offer, please visit <a href=\"http://www.xn-----6kcabdfroa7c7a2as1an7a2j.xn--p1ai/components/com_contact/views/categories/tmpl/5f9506d3f8.html\"><b>Our Site</b></a><br> <br> Best regards!<br><o:p></o:p></p></div></body></html> ------=_NextPart_000_0032_01D2912F.05324BC6--  "
			Ω(m).Should(Equal(
				Mail{
					Key:     "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
					Subject: &subject,
					Body:    &body,
					Junk:    true,
				}))
		})
	})
})
