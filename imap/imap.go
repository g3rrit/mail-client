package imap

import (
	"log"
	"io/ioutil"
	"strconv"

	"github.com/gproessl/mail-client/config"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
)

func RecvMail(cfg *config.Config) (error) {

	log.Println("Connecting to server...")

	c, err := client.DialTLS(cfg.ImapServer + ":" + cfg.ImapPort, nil)
	if err != nil {
		return err
	}
	defer c.Logout()
	log.Println("Connected")

	if err := c.Login(cfg.User, cfg.Pw); err != nil {
		return err
	}
	log.Println("Logged in")

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("Downloading " + m.Name)
		if err := downloadMailbox(cfg, c, m.Name); err != nil {
			return err
		}
	}

	if err := <-done; err != nil {
		return err
	}

	return nil
}

func downloadMailbox(cfg *config.Config, c *client.Client, mboxn string) (error) {
	mnums := make([]uint32, 256)

	files, err := ioutil.ReadDir(cfg.Maildir + mboxn)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() { continue }
		num, err := strconv.Atoi(file.Name())
		if err != nil {
			return err
		}
		mnums = append(mnums, uint32(num))
	}

	_, err = c.Select(mboxn, false)
	if err != nil {
		return err
	}

	log.Println("Fetching:", mnums)

	seqset := new(imap.SeqSet)
	seqset.AddNum(mnums...)
	items := []imap.FetchItem{imap.FetchAll}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()

	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
	}

	if err := <-done; err != nil {
		return err
	}

	return nil
}
