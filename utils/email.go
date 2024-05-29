package utils

import (
	"fmt"
	"net/smtp"
)

func SendMailSimple(subject string, body string, to []string) error {
	auth := smtp.PlainAuth(
		"",
		"akineshova00@gmail.com",
		"ghiauqstzujaalij",
		"smtp.gmail.com",
	)

	msg := "From: E-PLAYERS BLOG POST\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body

	err := smtp.SendMail(
		"smtp.gmail.com:587", // Note: Use port 587 for TLS
		auth,
		"akineshova00@gmail.com",
		to,
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
