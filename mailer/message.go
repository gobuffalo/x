package mailer

import (
	"io"

	"bytes"

	"github.com/gobuffalo/buffalo/render"
)

//Message represents an Email message
type Message struct {
	From    string
	To      []string
	CC      []string
	Bcc     []string
	Subject string

	Body        []byte
	ContentType string

	Attachments []Attachment
}

//Attachment are files added into a email message.
type Attachment struct {
	Name        string
	Reader      io.Reader
	ContentType string
}

// AddBody the message by receiving a renderer and rendering data.
func (m *Message) AddBody(r render.Renderer, data render.Data) error {
	buf := bytes.NewBuffer([]byte{})
	err := r.Render(buf, data)

	if err != nil {
		return err
	}

	m.Body = buf.Bytes()
	m.ContentType = r.ContentType()
	return nil
}

//AddAttachment adds the attachment to the list of attachments the Message has.
func (m *Message) AddAttachment(name, contentType string, r io.Reader) error {
	m.Attachments = append(m.Attachments, Attachment{
		Name:        name,
		ContentType: contentType,
		Reader:      r,
	})

	return nil
}

//NewMessage Builds a new message.
func NewMessage() Message {
	return Message{}
}
