package main

import (
	"log"
	"net/smtp"
)

func mail(to []string, message []byte) {
	sd := container.emailOpts

	from := sd.from
	pw := sd.password
	host := sd.mailHost
	addr := host + ":" + sd.mailPort
	message = fmtEmail(message, sd.template)

	auth := smtp.PlainAuth("", from, pw, host)

	err := smtp.SendMail(addr, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
}

func fmtEmail(message []byte, template []byte) []byte {

	panic("fmtEmail(): NOT IMPLEMENTED.")
}
