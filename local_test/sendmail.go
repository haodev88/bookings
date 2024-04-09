package main

import (
	"github.com/haodev88/bookings/internal/config"
	"github.com/haodev88/bookings/internal/helpers"
	"github.com/haodev88/bookings/internal/models"
	"log"
	"net/smtp"
)

var app config.AppConfig

func main()  {
	emailChan:=make(chan models.MailData)
	app.MailChan = emailChan
	defer close(emailChan)
	msg:= models.MailData{
		To: "john@do.com",
		From: "me@here.com",
		Subject: "Some subject",
		Content: "",
	}
	helpers.ListenForMail(msg)
	app.MailChan <- msg
}



func sendEmail()  {
	from:="haotv360@gmail.com"
	auth:=smtp.PlainAuth("", from, "", "localhost")
	err := smtp.SendMail("localhost:1025", auth, from, []string{"you@there.com"}, []byte("hello word"))
	if err != nil {
		log.Println(err)
	}
}
