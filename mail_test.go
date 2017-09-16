package sisyphus_test

import (
	"sort"

	s "github.com/carlostrub/sisyphus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mailBy func(m1, m2 *s.Mail) bool

func (by mailBy) Sort(mails []*s.Mail) {
	ms := &mailSorter{
		mails: mails,
		by:    by,
	}
	sort.Sort(ms)
}

type mailSorter struct {
	mails []*s.Mail
	by    func(m1, m2 *s.Mail) bool
}

func (ms *mailSorter) Len() int {
	return len(ms.mails)
}

func (ms *mailSorter) Swap(i, j int) {
	ms.mails[i], ms.mails[j] = ms.mails[j], ms.mails[i]
}

func (ms *mailSorter) Less(i, j int) bool {
	return ms.by(ms.mails[i], ms.mails[j])
}

var _ = Describe("Mail", func() {

	Context("Maildir", func() {
		It("Create a slice of mail keys", func() {
			result, err := s.Maildir("test/Maildir").Index()
			Ω(err).ShouldNot(HaveOccurred())

			name := func(m1, m2 *s.Mail) bool {
				return m1.Key < m2.Key
			}
			mailBy(name).Sort(result)
			Ω(result).Should(Equal(
				[]*s.Mail{

					{
						Key:  "1488181583.M633084P4781.mail.carlostrub.ch,S=708375,W=720014",
						Junk: true,
					},
					{
						Key:  "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
						Junk: true,
					},
					{
						Key:  "1488226337.M327824P8269.mail.carlostrub.ch,S=8044,W=8167",
						Junk: true,
					},
					{
						Key:  "1488226337.M327825P8269.mail.carlostrub.ch,S=802286,W=812785",
						Junk: true,
					},
					{
						Key:  "1488226337.M327833P8269.mail.carlostrub.ch,S=6960,W=7161",
						Junk: true,
					},
					{
						Key:  "1488228352.M339670P8269.mail.carlostrub.ch,S=12659,W=12782",
						Junk: true,
					},
					{
						Key: "1488230510.M141612P8565.mail.carlostrub.ch,S=5978,W=6119",
					},
					{
						Key:  "1504991721.M985788P1901.mail.carlostrub.ch,S=6474,W=6588",
						Junk: true,
					},
					{
						Key:  "1504991774.M467861P1924.mail.carlostrub.ch,S=6478,W=6592",
						Junk: true,
					},
					{
						Key:  "1505075914.M288773P9791.mail.carlostrub.ch,S=21241,W=21583",
						Junk: true,
					},
					{
						Key:  "1505392305.M710650P33881.mail.carlostrub.ch,S=6961,W=7064",
						Junk: true,
					},
				}))
		})
		It("Fail if Maildir does not exist", func() {
			_, err := s.Maildir("test/DOESNOTEXIST").Index()
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("Mail", func() {
		It("Load mail content into struct", func() {
			m := s.Mail{
				Key:     "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
				Subject: nil,
				Body:    nil,
				Junk:    true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			subject := "hello"
			body := "This is a multi-part message in MIME format.  ------=_NextPart_000_0032_01D2912F.05324BC6 Content-Type: text/plain; \tcharset=\"cp-850\" Content-Transfer-Encoding: quoted-printable  Dear cs,  We are looking for employees working remotely.  My name is Kari, I am the personnel manager of a large International company. Most of the work you can do from home, that is, at a distance.  Salary is $2000-$5300.  If you are interested in this offer, please visit  Our Site  Best regards! ------=_NextPart_000_0032_01D2912F.05324BC6 Content-Type: text/html; \tcharset=\"cp-850\" Content-Transfer-Encoding: quoted-printable  <html xmlns:v=\"urn:schemas-microsoft-com:vml\" xmlns:o=\"urn:schemas-microsoft-com:office:office\" xmlns:w=\"urn:schemas-microsoft-com:office:word\" xmlns:m=\"http://schemas.microsoft.com/office/2004/12/omml\" xmlns=\"http://www.w3.org/TR/REC-html40\"><head><META HTTP-EQUIV=\"Content-Type\" CONTENT=\"text/html; charset=us-ascii\"><meta name=Generator content=\"Microsoft Word 14 (filtered medium)\"><style><!-- /* Font Definitions */ @font-face \t{font-family:Calibri; \tpanose-1:2 15 5 2 2 2 4 3 2 4;} /* Style Definitions */ p.MsoNormal, li.MsoNormal, div.MsoNormal \t{margin:0in; \tmargin-bottom:.0001pt; \tfont-size:11.0pt; \tfont-family:\"Calibri\",\"sans-serif\";} a:link, span.MsoHyperlink \t{mso-style-priority:99; \tcolor:blue; \ttext-decoration:underline;} a:visited, span.MsoHyperlinkFollowed \t{mso-style-priority:99; \tcolor:purple; \ttext-decoration:underline;} span.EmailStyle17 \t{mso-style-type:personal-compose; \tfont-family:\"Calibri\",\"sans-serif\"; \tcolor:windowtext;} .MsoChpDefault \t{mso-style-type:export-only; \tfont-family:\"Calibri\",\"sans-serif\";} @page WordSection1 \t{size:8.5in 11.0in; \tmargin:1.0in 1.0in 1.0in 1.0in;} div.WordSection1 \t{page:WordSection1;} --></style><!--[if gte mso 9]><xml> <o:shapedefaults v:ext=\"edit\" spidmax=\"1026\" /> </xml><![endif]--><!--[if gte mso 9]><xml> <o:shapelayout v:ext=\"edit\"> <o:idmap v:ext=\"edit\" data=\"1\" /> </o:shapelayout></xml><![endif]--></head><body lang=EN-US link=blue vlink=purple><div class=WordSection1><p class=MsoNormal>Dear cs,<br> <br> We are looking for employees working remotely.<br> <br> My name is Kari, I am the personnel manager of a large International company.<br> Most of the work you can do from home, that is, at a distance.<br> <b>Salary is $2000-$5300.</b><br> <br> If you are interested in this offer, please visit <a href=\"http://www.xn-----6kcabdfroa7c7a2as1an7a2j.xn--p1ai/components/com_contact/views/categories/tmpl/5f9506d3f8.html\"><b>Our Site</b></a><br> <br> Best regards!<br><o:p></o:p></p></div></body></html> ------=_NextPart_000_0032_01D2912F.05324BC6--  "
			Ω(m).Should(Equal(
				s.Mail{
					Key:     "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
					Subject: &subject,
					Body:    &body,
					Junk:    true,
				}))
		})
		It("Fail if Subject has already content", func() {
			st := "test"
			m := s.Mail{
				Key:     "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
				Subject: &st,
				Body:    nil,
				Junk:    true,
			}

			err := m.Load("test/Maildir")
			Ω(err).Should(HaveOccurred())
		})
		It("Fail if Body has already content", func() {
			b := "test"
			m := s.Mail{
				Key:     "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
				Subject: nil,
				Body:    &b,
				Junk:    true,
			}

			err := m.Load("test/Maildir")
			Ω(err).Should(HaveOccurred())
		})

		It("Clean regular mail content", func() {
			m := s.Mail{
				Key:     "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
				Subject: nil,
				Body:    nil,
				Junk:    true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			subjectOutput := "hello"
			bodyOutput := " ------ 000 0032 01d2912f.05324bc6 : ; cp-850 : dear cs we are looking for employees working remotely my name is kari i am the personnel manager of a large international company most of the work you can do from home that is at a distance salary is 2000- 5300 if you are interested in this offer please visit our site best regards ------ 000 0032 01d2912f.05324bc6 : ; cp-850 : dear cs we are looking for employees working remotely. my name is kari i am the personnel manager of a large international company. most of the work you can do from home that is at a distance. salary is 2000- 5300. if you are interested in this offer please visit our site best regards ------ 000 0032 01d2912f.05324bc6-- "
			Ω(m).Should(Equal(
				s.Mail{
					Key:     "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730",
					Subject: &subjectOutput,
					Body:    &bodyOutput,
					Junk:    true,
				}))
		})

		It("Clean mail with base64 content", func() {
			m := s.Mail{
				Key:     "1488181583.M633084P4781.mail.carlostrub.ch,S=708375,W=720014:2,a",
				Subject: nil,
				Body:    nil,
				Junk:    true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			subjectOutput := "confirm remittance"
			bodyOutput := " ------ 000 0085 01c2a9a6.0d17aca6 : ; ---- 001 0086 01c2a9a6.0d17aca6 ------ 001 0086 01c2a9a6.0d17aca6 : ; : 7bit pfa remmittance copy value date 27022017 confirm payment detail thanks best regards admin director alliance bank this e-mail has been scanned for all known computer viruses this e-mail and any files transmitted with it are confidential and intended solely for the use of the individual or entity to whom they are addressed if you are not the intended recipient you are hereby notified that any dissemination forwarding copying or use of any of the information is strictly prohibited and the e-mail should immediately be deleted cobantur boltas makes no warranty as to the accuracy or completeness of any information contained in this message and hereby excludes any liability of any kind for the information contained therein or for the information transmission reception storage or use of such in any way whatsoever the opinions expressed in this message belong to sender alone and may not necessarily reflect the opinions of cobantur boltas ------ 001 0086 01c2a9a6.0d17aca6 : ; index.jpeg "
			Ω(m).Should(Equal(
				s.Mail{
					Key:     "1488181583.M633084P4781.mail.carlostrub.ch,S=708375,W=720014:2,a",
					Subject: &subjectOutput,
					Body:    &bodyOutput,
					Junk:    true,
				}))
		})

		It("More Junk", func() {
			m := s.Mail{
				Key:     "1488226337.M327824P8269.mail.carlostrub.ch,S=8044,W=8167:2,Sa",
				Subject: nil,
				Body:    nil,
				Junk:    true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			subjectOutput := "herpes breakthrough shocks medical world"
			bodyOutput := "--5ba77035ef2d5e8c615e79c26f9448f3 : ; : 8bit --5ba77035ef2d5e8c615e79c26f9448f3 : ; : 8bit i got herpes from this girl at a club but i got rid of it fast with this alert: herpes finally cured by rachael rettner senior writer february 27 2017 studies in mice suggest that gut bacteria can influence anxiety and other mental states. credit: dreamstime view full size image a new drug has successfully combated the virus that causes genital herpes starting today it will be used as a treatment for people with the condition. there have been many topical creams and drugs used as herpes cure treatments these treatments for herpes give short-term relief but only this can remove the virus and prevent re-occurrences to cure herpes. end your embarrassment - cure your herpes were appointed as provincial governors alongside members of the local aristocracy the title of doux was used but unlike earlier times these were mostly civilian governors with little military authority theodore awarded titles with such largesse that erly exclusive titles such as pansebastos sebastos or megalodoxotatos were devalued and came to be held by city notables to secure his new capital theodore instituted a guard of tzakones under a kastrophylax he portrait of a middle-aged man with a dark forked beard wearing a golden jewel-encrusted domed crown john iii doukas vatatzes emperor of nicaea from a 15th-century manuscript of the extracts of history of john zonaras --5ba77035ef2d5e8c615e79c26f9448f3-- "
			Ω(m).Should(Equal(
				s.Mail{
					Key:     "1488226337.M327824P8269.mail.carlostrub.ch,S=8044,W=8167:2,Sa",
					Subject: &subjectOutput,
					Body:    &bodyOutput,
					Junk:    true,
				}))
		})

		It("More Junk", func() {
			m := s.Mail{
				Key:     "1488226337.M327825P8269.mail.carlostrub.ch,S=802286,W=812785",
				Subject: nil,
				Body:    nil,
				Junk:    true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			subjectOutput := "cosan day 2017 new york friday march 24"
			bodyOutput := " --97bb499b431c6ea9472ef64c004b1106 : ; b82d2d2d1215b6505c0692de003f7694 --b82d2d2d1215b6505c0692de003f7694 : ; 075713cc6dc6562a324d22c8001fd90a --075713cc6dc6562a324d22c8001fd90a : ; : invitation cosan day 2017 new york friday march 24 2017 venue: park hyatt new york 153 west 57th street between 6th and 7th avenue new york ny 10019 the onyx room second level program 08:30 am registration 09:00 am cosan s/a csan3 presentations and q amp;a 10:45 am rumo s/a rumo3 presentation and q amp;a 11:30 am cosan limited czz presentation and q amp;a 12:10 pm closing and lunch rsvp http://www.invite-taylor-rafferty.com/ cosan/irday2017/default.htm or call briget ampudia at taylor rafferty 212 889 4350 or email cosan taylor-rafferty.com czz listed nyse csan3 novo mercado bm amp;fbovespa cgas5 cgas3 bm amp;fbovespa rlog3 novo mercado bm amp;fbovespa rumo3 novo mercado bm amp;fbovespa --075713cc6dc6562a324d22c8001fd90a : ; : if you cannot download the image below please view invitation and register online here --075713cc6dc6562a324d22c8001fd90a-- --b82d2d2d1215b6505c0692de003f7694 : ; cosan day 2017.jpg "
			Ω(m).Should(Equal(
				s.Mail{
					Key:     "1488226337.M327825P8269.mail.carlostrub.ch,S=802286,W=812785",
					Subject: &subjectOutput,
					Body:    &bodyOutput,
					Junk:    true,
				}))
		})

		It("More Junk", func() {
			m := s.Mail{
				Key:  "1488226337.M327833P8269.mail.carlostrub.ch,S=6960,W=7161:2,Sa",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			subjectOutput := "wear glasses your eyes are headed for serious trouble"
			bodyOutput := "--c2389532b48d1db204cfca8189242aeb : ; : 8bit --c2389532b48d1db204cfca8189242aeb : ; : 8bit snc .container width: 420px; .container .columns .container .column margin: 0; .container .fourteen.columns .container .fifteen.columns .container .sixteen.columns .container .one-third.column .container .two-thirds.column width: 420px; / self clearing goodness / .container:after content: 0020 ; display: block; height: 0; clear: both; visibility: hidden; .clearfix:before .clearfix:after .row:before .row:after content: 0020 ; display: block; overflow: hidden; visibility: hidden; width: 0; height: 0; .row:after .clearfix:after clear: both; .row .clearfix zoom: 1; .clear clear: both; display: block; overflow: hidden; visibility: hidden; width: 0; height: 0; if you wear glasses contacts or even if you think your vision can be improved you need to know about this.. in the link below you ll discover 1 weird trick that will drastically improve your vision gt; 1 trick to improve your vision today to your success 1 place ville marie 39th floor montreal quebec h3b4m7 canada email marketing by unsu bscribe westcott railway station served the village of westcott buckinghamshire near baron ferdinand de rothschild s estate at manor it was built by the duke of buckingham in 1871 as part of a short horse-drawn tramway that met the aylesbury and buckingham railway at quainton the next year it was converted for passenger use extended to brill railway station and renamed the brill tramway the poor quality locomotives running on the built and line were very slow initially limited to 5 miles per hour 8 km/h the line was taken over by the metropolitan railway in 1899 and transferred to public ownership in 1933 westcott station became part of the london underground despite being over 40 miles 60 km from central london until the closure of the line in 1935 the station building and its associated house pictured are the only significant buildings from the brill tramway to survive other than the junction station at quainton full article.. --c2389532b48d1db204cfca8189242aeb-- "
			Ω(m).Should(Equal(
				s.Mail{
					Key:     "1488226337.M327833P8269.mail.carlostrub.ch,S=6960,W=7161:2,Sa",
					Subject: &subjectOutput,
					Body:    &bodyOutput,
					Junk:    true,
				}))
		})

		It("More Junk", func() {
			m := s.Mail{
				Key:  "1488228352.M339670P8269.mail.carlostrub.ch,S=12659,W=12782:2,Sa",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			subjectOutput := "always in good form with our viagra super active."
			bodyOutput := " body .maintable height:100 important; width:100 important; margin:0; padding:0; img a img border:0; outline:none; text-decoration:none; .imagefix display:block; table td border-collapse:collapse; mso-table-lspace:0pt; mso-table-rspace:0pt; p margin:0; padding:0; margin-bottom:0; .readmsgbody width:100 ; .externalclass width:100 ; .externalclass .externalclass p .externalclass span .externalclass font .externalclass td .externalclass div line-height:100 ; img -ms-interpolation-mode: bicubic; body table td p a li blockquote -ms-text-size-adjust:100 ; -webkit-text-size-adjust:100 ; 96 \u00a0 if you can t read this email please view it online http://6url.ru/lhcj \u00a0 most popular products and special deals limited time offer hola the leading online store presents pharmaceuticals with delivery service in europe the united states and canada you can buy anti-acidity antifungals blood pressure herpes medication antifungals antibiotics anti-depressant diabetes medication antiviral anti-allergy/asthma and other various products keep your eye out for discount when purchasing\u00a0\u00a0\u00a0 check it now amazon web services inc is a subsidiary of amazon.com inc amazon.com is a registered trademark of amazon.com inc this message was produced and distributed by amazon web services inc 410 terry ave north seattle.https://aws.amazon.com/support if you no longer wish to receive these emails simply click on the following link unsubscribe © 2016 amazon all rights reserved \u00a0 "
			Ω(m).Should(Equal(
				s.Mail{
					Key:     "1488228352.M339670P8269.mail.carlostrub.ch,S=12659,W=12782:2,Sa",
					Subject: &subjectOutput,
					Body:    &bodyOutput,
					Junk:    true,
				}))
		})

		It("Wordlist 1", func() {
			m := s.Mail{
				Key:  "1488181583.M633084P4781.mail.carlostrub.ch,S=708375,W=720014:2,a",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"accuracy", "addressed", "admin", "alliance", "alone", "bank", "been", "belong", "best", "boltas", "cobantur", "computer", "confirm", "contained", "copy", "copying", "date", "deleted", "detail", "director", "entity", "excludes", "expressed", "files", "forwarding", "hereby", "individual", "intended", "kind", "known", "liability", "makes", "message", "notified", "opinions", "payment", "prohibited", "reception", "recipient", "reflect", "regards", "remittance", "scanned", "sender", "should", "solely", "storage", "strictly", "such", "thanks", "that", "therein", "they", "this", "value", "viruses", "warranty", "whatsoever", "whom", "with"}))
		})

		It("Wordlist 2", func() {
			m := s.Mail{
				Key:  "1488226337.M327822P8269.mail.carlostrub.ch,S=3620,W=3730:2,Sa",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"best", "company", "dear", "distance", "employees", "from", "hello", "home", "interested", "kari", "large", "looking", "manager", "most", "name", "offer", "personnel", "please", "regards", "remotely", "salary", "site", "that", "this", "visit", "work", "working"}))
		})

		It("Wordlist 3", func() {
			m := s.Mail{
				Key:  "1488226337.M327824P8269.mail.carlostrub.ch,S=8044,W=8167:2,Sa",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"alongside", "anxiety", "appointed", "authority", "awarded", "bacteria", "beard", "been", "came", "capital", "causes", "city", "civilian", "club", "combated", "creams", "crown", "cure", "cured", "dark", "devalued", "domed", "doukas", "doux", "dreamstime", "drug", "drugs", "earlier", "emperor", "erly", "exclusive", "extracts", "fast", "february", "finally", "forked", "from", "full", "genital", "girl", "give", "golden", "governors", "guard", "have", "held", "herpes", "history", "image", "influence", "instituted", "john", "largesse", "little", "local", "manuscript", "many", "medical", "members", "mental", "mice", "military", "mostly", "nicaea", "notables", "only", "other", "people", "portrait", "prevent", "provincial", "rachael", "relief", "remove", "rettner", "sebastos", "secure", "senior", "shocks", "size", "starting", "studies", "such", "suggest", "that", "theodore", "there", "these", "this", "times", "title", "titles", "today", "topical", "treatment", "treatments", "tzakones", "under", "unlike", "used", "vatatzes", "view", "virus", "wearing", "were", "will", "with", "world", "writer", "your", "zonaras"}))
		})

		It("Wordlist 4", func() {
			m := s.Mail{
				Key:  "1488226337.M327825P8269.mail.carlostrub.ch,S=802286,W=812785:2,Sa",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"ampudia", "avenue", "below", "between", "briget", "call", "cannot", "closing", "cosan", "download", "email", "friday", "here", "hyatt", "image", "invitation", "level", "limited", "listed", "lunch", "march", "mercado", "novo", "nyse", "online", "onyx", "park", "please", "program", "rafferty", "register", "room", "rsvp", "rumo", "second", "street", "taylor", "view", "west", "york"}))
		})

		It("Wordlist 5", func() {
			m := s.Mail{
				Key:  "1488226337.M327833P8269.mail.carlostrub.ch,S=6960,W=7161:2,Sa",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"about", "associated", "aylesbury", "baron", "became", "being", "below", "brill", "bscribe", "buckingham", "building", "buildings", "built", "canada", "central", "clearing", "closure", "contacts", "converted", "despite", "discover", "duke", "email", "estate", "even", "extended", "eyes", "ferdinand", "floor", "from", "full", "glasses", "goodness", "headed", "hour", "house", "improve", "improved", "initially", "junction", "know", "limited", "line", "link", "london", "manor", "marie", "marketing", "miles", "montreal", "near", "need", "next", "only", "other", "over", "ownership", "part", "passenger", "pictured", "place", "poor", "public", "quainton", "quality", "quebec", "railway", "renamed", "rothschild", "running", "self", "serious", "served", "short", "slow", "station", "success", "survive", "taken", "than", "that", "think", "today", "tramway", "trick", "trouble", "unsu", "until", "very", "village", "ville", "vision", "wear", "weird", "were", "westcott", "will", "year", "your"}))
		})

		It("Wordlist 6", func() {
			m := s.Mail{
				Key:  "1488228352.M339670P8269.mail.carlostrub.ch,S=12659,W=12782:2,Sa",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"always", "amazon", "antiviral", "blockquote", "blood", "body", "canada", "check", "click", "deals", "delivery", "diabetes", "discount", "email", "emails", "europe", "following", "font", "form", "good", "herpes", "hola", "keep", "leading", "limited", "link", "longer", "medication", "message", "most", "north", "offer", "online", "other", "please", "popular", "presents", "pressure", "produced", "products", "read", "receive", "registered", "reserved", "rights", "service", "services", "simply", "span", "special", "states", "store", "subsidiary", "super", "table", "terry", "these", "this", "time", "trademark", "united", "various", "viagra", "view", "when", "wish", "with", "your"}))
		})

		It("Wordlist 7", func() {
			m := s.Mail{
				Key: "1488230510.M141612P8565.mail.carlostrub.ch,S=5978,W=6119",
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"amending", "both", "build", "builds", "clang", "convert", "danfe", "depends", "drop", "explicit", "fine", "install", "instead", "library", "localbase", "manually", "port", "powerpc", "prefer", "rather", "shared", "static", "than", "their", "uses", "utilize", "with", "xorg"}))
		})

		It("Wordlist 8", func() {
			m := s.Mail{
				Key:  "1504991721.M985788P1901.mail.carlostrub.ch,S=6474,W=6588:2,S",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"‰", "⒏", "。", "《", "》", "下", "专", "倍", "六", "册", "利", "即", "取", "可", "合", "员", "回", "址", "够", "大", "天", "就", "彩", "拵", "拿", "提", "有", "永", "注", "澳", "特", "琻", "碼", "网", "赢", "邀", "钱", "门", "限", "領", "餸", "馈", "首", "，", "："}))
		})

		It("Wordlist 9", func() {
			m := s.Mail{
				Key:  "1504991774.M467861P1924.mail.carlostrub.ch,S=6478,W=6592:2,S",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"‰", "⒏", "。", "《", "》", "下", "专", "倍", "六", "册", "利", "即", "取", "可", "合", "员", "回", "址", "够", "大", "天", "就", "彩", "拵", "拿", "提", "有", "永", "注", "澳", "特", "琻", "碼", "网", "赢", "邀", "钱", "门", "限", "領", "餸", "馈", "首", "，", "："}))
		})

		It("Wordlist 10", func() {
			m := s.Mail{
				Key:  "1505392305.M710650P33881.mail.carlostrub.ch,S=6961,W=7064:2,S",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"agbetome", "banka", "drahy", "eddie", "fond", "odpov", "pozdravem", "prosim", "strycovy", "zesnuly"}))
		})

		It("Wordlist 11", func() {
			Skip("See known issues with mime/quotedprintable")
			m := s.Mail{
				Key:  "1505075914.M288773P9791.mail.carlostrub.ch,S=21241,W=21583:2,S",
				Junk: true,
			}

			err := m.Load("test/Maildir")
			Ω(err).ShouldNot(HaveOccurred())

			err = m.Clean()
			Ω(err).ShouldNot(HaveOccurred())

			list, err := m.Wordlist()
			Ω(err).ShouldNot(HaveOccurred())
			sort.Strings(list)

			Ω(list).Should(Equal(
				[]string{"agbetome", "banka", "drahy", "eddie", "fond", "odpov", "pozdravem", "prosim", "strycovy", "zesnuly"}))
		})
	})
})
