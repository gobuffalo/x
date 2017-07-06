package fakesmtp

import (
	"bufio"
	"net"
	"strings"
	"time"
)

type Server struct {
	Messages []string
	Listener net.Listener
}

//Start listens for connections on the given port
func (s *Server) Start(port string) error {

	conn, err := s.Listener.Accept()
	if err != nil {
		return err
	}

	s.Handle(&Connection{
		conn:    conn,
		address: conn.RemoteAddr().String(),
		time:    time.Now().Unix(),
		bufin:   bufio.NewReader(conn),
		bufout:  bufio.NewWriter(conn),
	})

	return nil
}

func (s *Server) Handle(c *Connection) {
	s.Messages = append(s.Messages, "")

	s.readHello(c)
	s.readSender(c)
	s.readRecipients(c)
	s.readData(c)

	c.conn.Close()
}

func (s *Server) readHello(c *Connection) {
	c.write("220 Welcome")
	text := c.read()
	s.AddMessageLine(text)

	c.write("250 Received")
}
func (s *Server) readSender(c *Connection) {
	text := c.read()
	s.AddMessageLine(text)
	c.write("250 Sender")
}

func (s *Server) readRecipients(c *Connection) {
	text := c.read()
	s.AddMessageLine(text)

	c.write("250 Recipient")
	text = c.read()
	for strings.Contains(text, "RCPT") {
		s.AddMessageLine(text)
		c.write("250 Recipient")
		text = c.read()
	}
}

func (s *Server) readData(c *Connection) {
	c.write("354 Ok Send data ending with <CRLF>.<CRLF>")

	for {
		text := c.read()
		bytes := []byte(text)
		s.AddMessageLine(text)
		// 46 13 10
		if bytes[0] == 46 && bytes[1] == 13 && bytes[2] == 10 {
			break
		}
	}
	c.write("250 server has transmitted the message")
}

func (s *Server) AddMessageLine(text string) {
	s.Messages[len(s.Messages)-1] = s.LastMessage() + text
}

func (s *Server) LastMessage() string {
	if len(s.Messages) == 0 {
		return ""
	}

	return s.Messages[len(s.Messages)-1]
}

func (s *Server) Clear() {
	s.Messages = []string{}
}

func NewServer(port string) (*Server, error) {
	s := &Server{Messages: []string{}}

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return s, err
	}
	s.Listener = listener
	return s, nil
}
