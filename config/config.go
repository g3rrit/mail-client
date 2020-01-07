package config

import (
	"encoding/json"
	"io/ioutil"
)

// Example Config:
// { 
// "ImapServer" : "www.imap.com",
// "ImapPort" : "993",
// "SmtpServer" : "www.smtp.com",
// "SmtpPort" : "993",
// "User" : "user",
// "Mail" : "user@mail.com",
// "Pw" : "pw",
// "Maildir" : "~/mail"
// }

type Config struct {
	ImapServer string
	ImapPort   string
	SmtpServer string
	SmtpPort   string
	User       string
	Mail       string
	Pw         string
	Maildir    string
}

func Load(path string) (Config, error) {
	config := Config{}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(dat, &config)
	return config, err
}
