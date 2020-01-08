package imap

import (
	"log"
	"io/ioutil"
	"strconv"
	"os"

	"github.com/gproessl/mail-client/config"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
)

func RecvMail(cfg *config.Config) (error) {

	log.Println("Connecting to server...")

	c, err := client.Dial(cfg.ImapServer + ":" + cfg.ImapPort)
	if err != nil {
		return err
	}
	defer c.Logout()
	log.Println("Connected")

	if err := c.StartTLS(nil); err != nil {
		return err
	}
	log.Println("TLS session initiated")

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
	// Selecting Mailbox
	mbox, err := c.Select(mboxn, false)
	if err != nil {
		return err
	}
	if mbox.Messages == 0 {
		log.Println("Zero Messages in", mbox.Name)
		return nil
	}

	// Creating local dir
	path := cfg.Maildir + "/" + mboxn
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("Creating dir:", path)
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// Scanning for local messages
	mnums := make([]uint32, 0)
	for _, file := range files {
		if file.IsDir() { continue }
		num, err := strconv.Atoi(file.Name())
		if err != nil {
			return err
		}
		mnums = append(mnums, uint32(num))
	}

	mgets := make([]uint32, 0)
	// TODO: refractor this mess
	for m := uint32(1); m <= mbox.Messages; m++ {
		if func(a uint32, b []uint32) bool {
			for _, c := range b { if a == uint32(c) { return false } }
			return true
		} (m, mnums) {
			mgets = append(mgets, m)
		}
	}

	log.Println("Fetching:", len(mgets), "of", mbox.Messages, "messages")

	if len(mgets) == 0 {
		return nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(mgets...)
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody, imap.FetchRFC822Text}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()

	for msg := range messages {
		log.Println("-*- " + msg.Envelope.Subject, msg.SeqNum)

		dat := make([]byte, 0)
		for _, v := range msg.Body {
			body := make([]byte, v.Len())
			_, err = v.Read(body)
			if err != nil {
				return err
			}
			dat = append(dat, body...)
		}

		if len(dat) == 0 {
			return nil
		}

		ioutil.WriteFile(path + "/" + strconv.Itoa(int(msg.SeqNum)), dat, 0755)
	}

	if err := <-done; err != nil {
		return err
	}

	return nil
}
