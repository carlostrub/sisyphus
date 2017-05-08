package sisyphus

import (
	"bufio"
	"errors"
	"log"
	"math"
	"mime/quotedprintable"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/kennygrant/sanitize"
	"github.com/luksen/maildir"
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

	log.Println("create missing directories for Maildir " + dir)

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

	log.Println("start indexing mails in " + dir)
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

	log.Println("all mails in " + dir + " indexed")

	return m, nil
}

// Load reads a mail's subject and body
func (m *Mail) Load(d string) error {

	if m.Junk {
		d = d + "/.Junk"
	}
	message, err := maildir.Dir(d).Message(m.Key)
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
func wordlist(s string) (l []string) {
	list := make(map[string]int)

	raw := strings.Split(s, " ")
	var clean []string

	// use regexp compile for use in the loop that follows
	regexMatcher, _ := regexp.Compile("(^[a-z]+$)")

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

	return l
}

// Wordlist prepares the mail for training
func (m *Mail) Wordlist() (w []string) {
	var s string

	if m.Subject != nil {
		s = s + " " + *m.Subject
	}

	if m.Body != nil {
		s = s + " " + *m.Body
	}

	w = wordlist(s)

	return w
}

// Classify analyses a new mail (a mail that arrived in the "new" directory),
// decides whether it is junk and -- if so -- moves it to the Junk folder. If
// it is not junk, the mail is untouched so it can be handled by the mail
// client.
func (m *Mail) Classify(db *bolt.DB) error {

	err := m.Clean()
	if err != nil {
		return err
	}

	list := m.Wordlist()
	junk, err := Junk(db, list)
	if err != nil {
		return err
	}

	log.Print("Classified " + m.Key + " as Junk=" + strconv.FormatBool(m.Junk))

	// Move mail around if junk.
	if junk {
		m.Junk = junk
		err := os.Rename("./new/"+m.Key, "./.Junk/cur/"+m.Key)
		if err != nil {
			return err
		}
		log.Print("Moved " + m.Key + " from new to Junk folder")
	}

	return nil
}

// Learn adds the words to the respective list and unlearns on the other, if
// the mail has been moved from there.
func (m *Mail) Learn(db *bolt.DB) error {
	return nil
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
