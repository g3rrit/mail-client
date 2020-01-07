package main

import (
	"log"
	"os"

	"github.com/gproessl/mail-client/config"
	"github.com/gproessl/mail-client/imap"
	"github.com/gproessl/mail-client/smtp"
)

func printUsage() {
	log.Println("usage: mail-client config_file (s|r) [mail]")
}

func sendMail(cfg *config.Config, path string) {
	err := smtp.SendMail(cfg, path)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Sent Mail")
}

func recvMail(cfg *config.Config) {
	err := imap.RecvMail(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Received Mail")
}

func main() {
	if len(os.Args) < 3 {
		printUsage()
		return
	}

	cfg, err := config.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}

	switch os.Args[2] {
	case "s":
		if len(os.Args) < 4 {
			log.Fatal("mail to send not specified")
			return
		}
		log.Println("Sending Mail")
		sendMail(&cfg, os.Args[3])
	case "r":
		log.Println("Receiving Mail")
		recvMail(&cfg)
	default:
		printUsage()
	}
}
