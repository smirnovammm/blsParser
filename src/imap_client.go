package src

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
)

type ImapClient struct {
	client *client.Client
}

func NewImapClient(email string, password string) (*ImapClient, error) {
	log.Println("Connecting to mailbox")
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return nil, err
	}
	log.Println("Connected")

	log.Printf("Login to %s", email)
	if err := c.Login(email, password); err != nil {
		return nil, err
	}
	return &ImapClient{client: c}, nil
}

func (c ImapClient) GetLastOtp(afterTime *time.Time) (*string, *time.Time, error) {
	var code *string
	var date *time.Time
	var err error

	for i := 1; i < 4; i++ {
		code, date, err = c.getLastOtp()
		if err != nil {
			continue
		}
		if date.After(*afterTime) {
			return code, date, nil
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, fmt.Errorf("Failed to find otp after: %s Last received message at: %s ", afterTime, *date)
}

func (c ImapClient) getLastOtp() (*string, *time.Time, error) {
	var code string
	var date time.Time

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.client.List("", "*", mailboxes)
	}()

	if err := <-done; err != nil {
		return nil, nil, err
	}

	mbox, err := c.client.Select("bls-portugal", false)
	if err != nil {
		return nil, nil, err
	}

	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 1 {
		from = mbox.Messages - 1
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.client.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822}, messages)
	}()

	for msg := range messages {
		for _, r := range msg.Body {
			entity, err := message.Read(r)
			if message.IsUnknownCharset(err) {
			} else if err != nil {
				return nil, nil, err
			}

			bt, rErr := ioutil.ReadAll(entity.Body)
			if rErr != nil {
				return nil, nil, err
			}

			t := strings.Split(strings.Split(toUtf8(bt), ">Date : ")[1], "</td>")[0]
			log.Println(t)
			date, err = time.Parse("2006-01-02 15:04:05", t)
			if err != nil {
				return nil, nil, err
			}
			code = strings.Split(strings.Split(toUtf8(bt), "Verification code - ")[1], "<br/>")[0]
		}
	}

	return &code, &date, nil
}

func (c ImapClient) Logout() {
	if err := c.client.Logout(); err != nil {
		panic(err)
	}
}

func toUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}
