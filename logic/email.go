package logic

import (
	"dunakeke/config"

	"gopkg.in/gomail.v2"
)

type Messagee struct {
    Name    string
    Email   string
}

var dialer *gomail.Dialer

func SetupEmail() {
    dialer = gomail.NewDialer(config.Config.Email.Smtp,  config.Config.Email.Port,
                              config.Config.Email.Email, config.Config.Email.Password)
}

func SendEmail(subject string, html string, to ...Messagee) {
    sender, err := dialer.Dial()
    if nil != err {
        log.Println("Faild to get email sender: ", err)
    }

    message := gomail.NewMessage()
    for _, t := range(to) {
        message.SetAddressHeader("From", config.Config.Email.Email, config.Config.Email.Sender)
        message.SetAddressHeader("To", t.Email, t.Name)
        message.SetHeader("Subject", subject)
        message.SetBody("text/html", html)

        err := gomail.Send(sender, message)
        if nil != err {
            log.Println("Unable to send message to: ", t)
        }

        message.Reset()
    }
}
