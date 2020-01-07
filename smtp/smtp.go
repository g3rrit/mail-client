package smtp

import (
	"io/ioutil"
	"strings"
	"errors"
	"log"
	"fmt"

	"net/smtp"

	"github.com/gproessl/mail-client/config"
)

func SendMail(cfg *config.Config, path string) (error) {
	mail, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	msg := string(mail)

	lines := strings.Split(msg, "\n")

	if lines[0][0:4] != "To: " {
		return errors.New("No recipient specified")
	}
	to := lines[0][4:]

	msg = strings.Join(lines, "\r\n")


	log.Println("Recipient:", to)
	log.Println("Content:\n", msg)

	fmt.Print("Send? (y/n): ")
	var res string
	fmt.Scanln(&res)
	if res != "y" {
		return errors.New("Not sending")
	}

	auth := smtp.PlainAuth("", cfg.User, cfg.Pw, cfg.SmtpServer)
	err = smtp.SendMail(cfg.SmtpServer + ":" + cfg.SmtpPort, auth, cfg.Mail, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
