package mailer_test

import (
	"testing"

	"github.com/gobuffalo/buffalo/render"
	server "github.com/gobuffalo/x/fakesmtp"
	"github.com/gobuffalo/x/mailer"
	"github.com/stretchr/testify/require"
)

var sender mailer.Deliverer
var rend *render.Engine

const smtpPort = "9807"

func init() {
	rend = render.New(render.Options{})
}

func TestSendPlain(t *testing.T) {
	server.WaitForMessage(smtpPort)

	r := require.New(t)
	smtp, err := mailer.NewSMTPMailer(smtpPort, "127.0.0.1", "username", "password")
	r.Nil(err)

	m := mailer.NewMessage()
	m.From = "mark@example.com"
	m.To = []string{"something@something.com"}
	m.Subject = "Cool Message"
	m.CC = []string{"other@other.com", "my@other.com"}
	m.Bcc = []string{"secret@other.com"}
	m.AddBody(rend.String("Hello <%= Name %>"), render.Data{"Name": "Antonio"}, "text/plain")
	r.Equal(m.Body, []byte("Hello Antonio"))

	err = smtp.Deliver(m)

	r.Contains(server.LastMessage, "FROM:<mark@example.com>")
	r.Contains(server.LastMessage, "RCPT TO:<other@other.com>")
	r.Contains(server.LastMessage, "RCPT TO:<my@other.com>")
	r.Contains(server.LastMessage, "RCPT TO:<secret@other.com>")
	r.Contains(server.LastMessage, "Subject: Cool Message")
	r.Contains(server.LastMessage, "Cc: other@other.com, my@other.com")
	r.Contains(server.LastMessage, "Content-Type: text/plain")
	r.Contains(server.LastMessage, "Hello Antonio")
}
