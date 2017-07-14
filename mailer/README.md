# github.com/gobuffalo/x/mailer

This package is intended to allow easy Email sending with Buffalo, it allows you to define your custom `mailer.Sender` for the provider you would like to use.

The following is an example on how to use this package with your Buffalo app:

```go
//actions/mailer.go

import (
    "github.com/gobuffalo/x/mailer"
    "github.com/gobuffalo/buffalo/render"ç
    "github.com/gobuffalo/packr"
    "github.com/me/myapp/models"
    "github.com/pkg/errors"
)

var smtp mailer.Sender

func init() {
    port := envy.Get("SMTP_PORT", "1025")
    host := envy.Get("SMTP_HOST", "localhost")
    user := envy.Get("SMTP_USER", "")
    password := envy.Get("SMTP_PASSWORD", "")

	var err error
	smtp, err = mailer.NewSMTPMailer(host,port, user, password)
	
    if err != nil {
		log.Fatal(err)
	}

	r = render.New(render.Options{
		TemplatesBox: packr.NewBox("../templates"),
	})
}

//SendContactMessage Sends contact message to contact@myapp.com
func SendContactMessage(c *models.Contact) error {
	m := mailer.NewMessage()
	m.Subject = "New Contact"
    m.To = []string {"contact@myapp.com"}

	data := map[string]interface{}{
		"contact": c,
	}
	
    err := m.AddBodies(data, r.HTML("mail/contact.md"), r.Plain("mail/contact.md"))

	if err != nil {
		return errors.WithStack(err)
	}

    err := smtp.Send(m)
}

```
