package main

import (
	"encoding/json"
	"net/smtp"
	"strings"
)

// TODO: Proper error handling

func Mail(message []byte) bool {
	sd := container.emailOpts

	from := sd.from
	pw := sd.password
	host := sd.mailHost
	addr := host + ":" + sd.mailPort

	req, err := parseMessage(message)
	if err != nil {
		return false
	}

	to := []string{req.to}
	message = []byte(fmtEmail(req, string(sd.template)))

	auth := smtp.PlainAuth("", from, pw, host)

	err = smtp.SendMail(addr, auth, from, to, message)
	if err != nil {
		return false
	}

	return true
}

func fmtEmail(request *RequestParams, template string) string {
	out := template
	out = strings.Replace(out, "{name}", request.name, 1)
	out = strings.Replace(out, "{url}", request.url, 1)

	return out
}

type RequestParams struct {
	to   string
	name string
	url  string
}

func parseMessage(message []byte) (*RequestParams, error) {
	var out RequestParams
	err := json.Unmarshal(message, &out)
	return &out, err
}
