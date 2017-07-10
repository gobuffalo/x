package mailer_test

import (
	"bytes"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/x/fakesmtp"
	"github.com/gobuffalo/x/mailer"
	"github.com/stretchr/testify/require"
)

var sender mailer.Deliverer
var rend *render.Engine
var smtpServer *fakesmtp.Server

const smtpPort = "2002"

func init() {
	rend = render.New(render.Options{})
	smtpServer, _ = fakesmtp.NewServer(smtpPort)
	sender, _ = mailer.NewSMTPMailer("127.0.0.1", smtpPort, "username", "password")

	go smtpServer.Start(smtpPort)
}

func TestSendPlain(t *testing.T) {
	smtpServer.Clear()
	r := require.New(t)

	m := mailer.Message{
		From:    "mark@example.com",
		To:      []string{"something@something.com"},
		Subject: "Cool Message",
		CC:      []string{"other@other.com", "my@other.com"},
		Bcc:     []string{"secret@other.com"},
	}

	m.AddAttachment("someFile.txt", "text/plain", bytes.NewBuffer([]byte("hello")))
	m.AddBody(rend.String("Hello <%= Name %>"), render.Data{"Name": "Antonio"})
	r.Equal(m.Bodies[0].Content, "Hello Antonio")

	err := sender.Deliver(m)
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
