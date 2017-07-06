package fakesmtp

import (
	"bufio"
	"fmt"
	"net"
)

type Connection struct {
	conn    net.Conn
	address string
	time    int64
	bufin   *bufio.Reader
	bufout  *bufio.Writer
}

func (c *Connection) write(s string) {
	c.bufout.WriteString(s + "\r\n")
	c.bufout.Flush()
}
func (c *Connection) read() string {
	reply, err := c.bufin.ReadString('\n')

	if err != nil {
		fmt.Println("e ", err)
	}
	return reply
}
