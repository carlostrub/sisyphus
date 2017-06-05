package sisyphus

import (
	"bufio"
	"errors"
	"math"
	"mime/quotedprintable"
	"net/mail"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/carlostrub/maildir"
	"github.com/kennygrant/sanitize"
)

// Maildir represents the address to a Maildir directory
type Maildir string

// Mail includes the key of a mail in Maildir
type Mail struct {
	Key           string
	Subject, Body *string
	Junk, New     bool
}

// CreateDirs creates all the required dirs -- if not already there.
func (d Maildir) CreateDirs() error {

	dir := string(d)

	log.WithFields(log.Fields{
		"dir": dir,
	}).Info("Create missing directories")

	err := os.MkdirAll(dir+"/.Junk/cur", 0700)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir+"/new", 0700)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir+"/cur", 0700)

	return err
}

// Index loads all mail keys from the Maildir directory for processing.
func (d Maildir) Index() (m []*Mail, err error) {

	dir := string(d)

	log.WithFields(log.Fields{
		"dir": dir,
	}).Info("Start indexing mails")

	dirs := []string{dir, dir + "/.Junk"}
	for _, val := range dirs {
		j, err := maildir.Dir(val).Keys()
		if err != nil {
			return m, err
		}
		for _, v := range j {
			var new Mail
			new.Key = v
			if val == dir+"/.Junk" {
				new.Junk = true
			}
			m = append(m, &new)
		}
	}

	log.WithFields(log.Fields{
		"dir": dir,
	}).Info("All mails indexed")

	return m, nil
}

// Load reads a mail's subject and body
func (m *Mail) Load(dir Maildir) (err error) {

	var message *mail.Message
	message = new(mail.Message)

	switch {
	case m.Junk:
		dir = dir + Maildir("/.Junk")
	case m.New:
		dir = dir + Maildir("/new")
	}

	message, err = maildir.Dir(dir).Message(m.Key)
	if err != nil {
		return err
	}

	// get Subject
	if m.Subject != nil {
		return errors.New("there is already a subject")
	}
	subject := message.Header.Get("Subject")
	m.Subject = &subject

	// get Body
	bQ := quotedprintable.NewReader(message.Body)
	var b []string
	bScanner := bufio.NewScanner(bQ)
	for bScanner.Scan() {
		raw := bScanner.Text()
		b = append(b, raw)
	}

	body := strings.Join(b, " ")
	if m.Body != nil {
		return errors.New("there is already a body")
	}
	m.Body = &body

	return nil
}

func trimStringFromBase64(s string) string {
	if idx := strings.Index(s, "Content-Transfer-Encoding: base64"); idx != -1 {
		return s[:idx-1]
	}
	return s
}

func cleanString(i string) (s string) {

	s = sanitize.Accents(i)
	s = sanitize.HTML(s)
	s = strings.ToLower(s)

	bad := []string{
		"boundary=", "charset", "content-transfer-encoding",
		"content-type", "image/jpeg", "multipart/alternative",
		"multipart/related", "name=", "nextpart", "quoted-printable",
		"text/html", "text/plain", "this email must be viewed in html mode",
		"this is a multi-part message in mime format",
		"windows-1251", "windows-1252", "!", "#", "$", "%", "&", "'",
		"(", ")", "*", "+", ",", ". ", "<", "=", ">", "?", "@", "[",
		"\"", "\\", "\n", "\t", "]", "^", "_", "{", "|", "}",
	}
	for _, b := range bad {
		s = strings.Replace(s, b, " ", -1)
	}
	for i := 0; i < 10; i++ {
		s = strings.Replace(s, "  ", " ", -1)
	}

	return s
}

// Clean cleans the mail's subject and body
func (m *Mail) Clean() error {
	if m.Subject != nil {
		s := trimStringFromBase64(*m.Subject)
		s = cleanString(s)
		m.Subject = &s
	}

	if m.Body != nil {
		b := trimStringFromBase64(*m.Body)
		b = cleanString(b)
		m.Body = &b
	}
	return nil
}

// wordlist takes a string of space separated text and returns a list of unique
// words in a space separated string
func wordlist(s string) (l []string, err error) {
	list := make(map[string]int)

	raw := strings.Split(s, " ")
	var clean []string

	// use regexp compile for use in the loop that follows
	regexMatcher, err := regexp.Compile("(^[a-z]+$)")
	if err != nil {
		return l, err
	}

	for _, w := range raw {

		// no long or too short words
		length := len(w)
		if length < 4 || length > 10 {
			continue
		}

		// no numbers, special characters, etc. -- only words
		match := regexMatcher.MatchString(w)
		if !match {
			continue
		} else {
			clean = append(clean, w)
		}
	}

	// only the first 200 words count
	maxWords := int(math.Min(200, float64(len(clean))))
	for i := 0; i < maxWords; i++ {
		w := clean[i]
		list[w]++
	}

	for word, count := range list {
		if count > 10 {
			continue
		}

		l = append(l, word)
	}

	return l, nil
}

// Wordlist prepares the mail for training
func (m *Mail) Wordlist() (w []string, err error) {
	var s string

	if m.Subject != nil {
		s = s + " " + *m.Subject
	}

	if m.Body != nil {
		s = s + " " + *m.Body
	}

	w, err = wordlist(s)

	return w, err
}

// cleanWordlist combines Clean and Wordlist in one internal function
func (m *Mail) cleanWordlist() (w []string, err error) {
	err = m.Clean()
	if err != nil {
		return w, err
	}

	w, err = m.Wordlist()

	return w, err
}

// LoadMails creates missing directories and then loads all mails from a given
// slice of Maildirs
func LoadMails(d []Maildir) (mails map[Maildir][]*Mail, err error) {
	mails = make(map[Maildir][]*Mail)

	// create missing directories and write index
	for _, val := range d {
		err := val.CreateDirs()
		if err != nil {
			return mails, err
		}

		var m []*Mail
		m, err = val.Index()
		if err != nil {
			return mails, err
		}

		mails[val] = m
	}

	return mails, nil
}
