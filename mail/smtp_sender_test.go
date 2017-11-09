package mail_test

import (
	"bytes"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/x/fakesmtp"
	"github.com/gobuffalo/x/mail"
	"github.com/stretchr/testify/require"
)

var sender mail.Sender
var rend *render.Engine
var smtpServer *fakesmtp.Server

const smtpPort = "2002"

func init() {
	rend = render.New(render.Options{})
	smtpServer, _ = fakesmtp.New(smtpPort)
	sender, _ = mail.NewSMTPSender("127.0.0.1", smtpPort, "username", "password")

	go smtpServer.Start(smtpPort)
}

func TestSendPlain(t *testing.T) {
	smtpServer.Clear()
	r := require.New(t)

	m := mail.NewMessage()
	m.From = "mark@example.com"
	m.To = []string{"something@something.com"}
	m.Subject = "Cool Message"
	m.CC = []string{"other@other.com", "my@other.com"}
	m.Bcc = []string{"secret@other.com"}

	m.AddAttachment("someFile.txt", "text/plain", bytes.NewBuffer([]byte("hello")))
	m.AddBody(rend.String("Hello <%= Name %>"), render.Data{"Name": "Antonio"})
	r.Equal(m.Bodies[0].Content, "Hello Antonio")

	err := sender.Send(m)
	r.Nil(err)

	lastMessage := smtpServer.LastMessage()

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
