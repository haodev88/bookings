package helpers

import (
	"github.com/haodev88/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"time"
)

func ListenForMail(m models.MailData)  {
	go sendMsg(m)
}

func sendMsg(m models.MailData)  {
	log.Println("Hao tran")
	server:= mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client,err:= server.Connect()
	if err != nil {
		log.Println(err)
	}

	email:= mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, "Hello, <strong>Word</strong>")
	// email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}
}