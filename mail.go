package main

import (
	"bufio"
	"errors"
	"mime/quotedprintable"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/luksen/maildir"
)

// Mail includes the key of a mail in Maildir
type Mail struct {
	Key           string
	Subject, Body *string
	Junk          bool
}

// Index loads all mail keys from the Maildir directory for processing.
func Index(d string) (m []*Mail, err error) {

	g, err := maildir.Dir(d).Keys()
	if err != nil {
		return m, err
	}
	for _, val := range g {
		var new Mail
		new.Key = val
		m = append(m, &new)
	}

	j, err := maildir.Dir(d + "/.Junk").Keys()
	if err != nil {
		return m, err
	}
	for _, val := range j {
		var new Mail
		new.Key = val
		new.Junk = true
		m = append(m, &new)
	}

	return m, nil
}

func trimStringFromBase64(s string) string {
	if idx := strings.Index(s, "Content-Transfer-Encoding: base64"); idx != -1 {
		return s[:idx-1]
	}
	return s
}

func cleanString(i string) (s string, err error) {

	s = trimStringFromBase64(i)
	s = sanitize.Accents(s)
	s = sanitize.HTML(s)
	s = strings.ToLower(s)

	s = strings.Replace(s, "boundary=", " ", -1)
	s = strings.Replace(s, "charset", " ", -1)
	s = strings.Replace(s, "content-transfer-encoding", " ", -1)
	s = strings.Replace(s, "content-type", " ", -1)
	s = strings.Replace(s, "cp-850", " ", -1)
	s = strings.Replace(s, "image/jpeg", " ", -1)
	s = strings.Replace(s, "multipart/alternative", " ", -1)
	s = strings.Replace(s, "multipart/related", " ", -1)
	s = strings.Replace(s, "name=", " ", -1)
	s = strings.Replace(s, "nextpart", " ", -1)
	s = strings.Replace(s, "quoted-printable", " ", -1)
	s = strings.Replace(s, "text/html", " ", -1)
	s = strings.Replace(s, "text/plain", " ", -1)
	s = strings.Replace(s, "this email must be viewed in html mode", " ", -1)
	s = strings.Replace(s, "this is a multi-part message in mime format", " ", -1)
	s = strings.Replace(s, "windows-1251", " ", -1)
	s = strings.Replace(s, "windows-1252", " ", -1)

	s = strings.Replace(s, "!", " ", -1)
	s = strings.Replace(s, "#", " ", -1)
	s = strings.Replace(s, "$", " ", -1)
	s = strings.Replace(s, "%", " ", -1)
	s = strings.Replace(s, "&", " ", -1)
	s = strings.Replace(s, "'", "", -1)
	s = strings.Replace(s, "(", " ", -1)
	s = strings.Replace(s, ")", " ", -1)
	s = strings.Replace(s, "*", " ", -1)
	s = strings.Replace(s, "+", " ", -1)
	s = strings.Replace(s, ",", " ", -1)
	s = strings.Replace(s, ". ", " ", -1)
	s = strings.Replace(s, "<", " ", -1)
	s = strings.Replace(s, "=", " ", -1)
	s = strings.Replace(s, ">", " ", -1)
	s = strings.Replace(s, "?", " ", -1)
	s = strings.Replace(s, "@", " ", -1)
	s = strings.Replace(s, "[", " ", -1)
	s = strings.Replace(s, "\"", " ", -1)
	s = strings.Replace(s, "\\", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)
	s = strings.Replace(s, "\t", " ", -1)
	s = strings.Replace(s, "]", " ", -1)
	s = strings.Replace(s, "^", " ", -1)
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Replace(s, "{", " ", -1)
	s = strings.Replace(s, "|", " ", -1)
	s = strings.Replace(s, "}", " ", -1)

	for i := 0; i < 10; i++ {
		s = strings.Replace(s, "  ", " ", -1)
	}

	return s, nil
}

// Clean prepares the mail's subject and body for training
func (m *Mail) Clean() error {
	if m.Subject != nil {
		s, err := cleanString(*m.Subject)
		if err != nil {
			return err
		}
		m.Subject = &s
	}

	if m.Body != nil {
		b, err := cleanString(*m.Body)
		if err != nil {
			return err
		}
		m.Body = &b
	}
	return nil
}

// Load reads a mail's subject and body
func (m *Mail) Load(d string) error {

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
