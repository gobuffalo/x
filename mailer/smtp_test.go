package mailer_test

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

var listening bool
var LastMessage string

type Client struct {
	conn    net.Conn
	address string
	time    int64
	bufin   *bufio.Reader
	bufout  *bufio.Writer
}

func (c *Client) write(s string) {
	c.bufout.WriteString(s + "\r\n")
	c.bufout.Flush()
}
func (c *Client) read() string {
	reply, err := c.bufin.ReadString('\n')

	if err != nil {
		fmt.Println("e ", err)
	}
	return reply
}

func appendToFile(text string) {
	LastMessage = LastMessage + text
}

func handleClient(c *Client) {
	LastMessage = ""
	c.write("220 Welcome to the Jungle")
	text := c.read()
	appendToFile(text)
	c.write("250 No one says helo anymore")
	text = c.read()
	appendToFile(text)
	c.write("250 Sender")
	text = c.read()
	appendToFile(text)

	c.write("250 Recipient")
	text = c.read()
	for strings.Contains(text, "RCPT") {
		appendToFile(text)
		c.write("250 Recipient")
		text = c.read()
	}

	c.write("354 Ok Send data ending with <CRLF>.<CRLF>")

	for {
		text = c.read()
		bytes := []byte(text)
		appendToFile(text)
		// 46 13 10
		if bytes[0] == 46 && bytes[1] == 13 && bytes[2] == 10 {
			break
		}
	}
	c.write("250 server has transmitted the message")
	c.conn.Close()
}

func StartSMTPServer(port string) {
	listening = true
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:"+port)
		if err != nil {
			fmt.Println("run as root")
			return
		}

		for listening {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleClient(&Client{
				conn:    conn,
				address: conn.RemoteAddr().String(),
				time:    time.Now().Unix(),
				bufin:   bufio.NewReader(conn),
				bufout:  bufio.NewWriter(conn),
			})
		}
	}()
}

func StopSMTPServer() {
	listening = false
}
