# github.com/gobuffalo/x/mailer

This package is intended to allow easy Email sending with Buffalo, it allows you to define your custom `mailer.Sender` for the provider you would like to use.

The following is an example on how to setup a smtp mailer as well as creating a Mailer function to use the 

```go
//actions/mailer.go

import (
    "github.com/gobuffalo/x/mailer"
    "github.com/gobuffalo/packr"
    "github.com/me/myapp/models"
    "github.com/pkg/errors"
	"github.com/gobuffalo/buffalo/render"
)

var smtp mailer.Sender
var r *render.Engine

func init() {
	port := envy.Get("SMTP_PORT", "1025")
	host := envy.Get("SMTP_HOST", "localhost")
	user := envy.Get("SMTP_USER", "")
	password := envy.Get("SMTP_PASSWORD", "")

	var err error
	smtp, err = mailer.NewSMTPSender(host, port, user, password)

	if err != nil {
		log.Fatal(err)
	}

	r = render.New(render.Options{
		HTMLLayout:     "application.html",
		TemplateEngine: plush.BuffaloRenderer,
		TemplatesBox:   packr.NewBox("../templates"),
		Helpers: map[string]interface{}{},
	}
}

//SendContactMessage Sends contact message to contact@myapp.com
func SendContactMessage(c *models.Contact) error {
    
	//Creates a new message
	m := mailer.NewMessage()
	m.From = "sender@myapp.com"
	m.Subject = "New Contact"
	m.To = []string{"contact@myapp.com"}

	// Data that will be used inside the templates when rendering.
	data := map[string]interface{}{
		"contact": c,
	}

	// You can add multiple bodies to the message you're creating to have content-types alternatives.
	err := m.AddBodies(data, r.HTML("mail/contact.html"), r.Plain("mail/contact.txt"))

	if err != nil {
		return errors.WithStack(err)
	}

	err = smtp.Send(m)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

```

This `SendContactMessage` could be called by one of your actions, p.e. the action that handles your contact form submission.

```go
//actions/contact.go
...

func ContactFormHandler(c buffalo.Context) error {
    contact := &models.Contact{}
    c.Bind(contact)
    
    //Calling to send the message
    SendContactMessage(contact)
    return c.Redirect(302, "contact/thanks")
}
...
```

If you're using Gmail or need to configure your SMTP connection you can use the Dialer property on the SMTPSender, p.e: (for Gmail)

```go
...
var smtp mailer.Sender

func init() {
    port := envy.Get("SMTP_PORT", "465") 
    // or 587 with TLS 

	host := envy.Get("SMTP_HOST", "smtp.gmail.com")
	user := envy.Get("SMTP_USER", "your@email.com")
	password := envy.Get("SMTP_PASSWORD", "yourp4ssw0rd")

	var err error
	sender, err := mailer.NewSMTPSender(host, port, user, password)
	sender.Dialer.SSL = true

    //or if TLS
    sender.Dialer.TLSConfig = &tls.Config{...}
    
    smtp = sender
}
...
```