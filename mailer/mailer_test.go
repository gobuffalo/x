package mailer_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/x/fakesmtp"
	"github.com/gobuffalo/x/mailer"
	"github.com/stretchr/testify/require"
)

var sender mailer.Deliverer
var rend *render.Engine
var smtpServer *fakesmtp.Server

const smtpPort = "9807"

func init() {
	rend = render.New(render.Options{})
}

func TestSendPlain(t *testing.T) {
	r := require.New(t)

	smtpServer, err := fakesmtp.NewServer(smtpPort)
	go smtpServer.Start(smtpPort)

	r.Nil(err)
	smtp, err := mailer.NewSMTPMailer(smtpPort, "127.0.0.1", "username", "password")
	r.Nil(err)

	m := mailer.NewMessage()
	m.From = "mark@example.com"
	m.To = []string{"something@something.com"}
	m.Subject = "Cool Message"
	m.CC = []string{"other@other.com", "my@other.com"}
	m.Bcc = []string{"secret@other.com"}
	m.AddAttachment("someFile.txt", "text/plain", bytes.NewBuffer([]byte("hello")))
	m.AddBody(rend.String("Hello <%= Name %>"), render.Data{"Name": "Antonio"}, "text/plain")
	r.Equal(m.Body, []byte("Hello Antonio"))

	err = smtp.Deliver(m)
	lastMessage := smtpServer.LastMessage()

	log.Println(smtpServer.Messages)

	r.Contains(lastMessage, "FROM:<mark@example.com>")
	r.Contains(lastMessage, "RCPT TO:<other@other.com>")
	r.Contains(lastMessage, "RCPT TO:<my@other.com>")
	r.Contains(lastMessage, "RCPT TO:<secret@other.com>")
	r.Contains(lastMessage, "Subject: Cool Message")
	r.Contains(lastMessage, "Cc: other@other.com, my@other.com")
	r.Contains(lastMessage, "Content-Type: text/plain")
	r.Contains(lastMessage, "Hello Antonio")
	r.Contains(lastMessage, "Content-Disposition: attachment; filename=\"someFile.txt\"")
	r.Contains(lastMessage, "aGVsbG8=") //base64 of the file content
}
