package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) *TelnetClient {
	// Place your code here.
	return &TelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *TelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}

	c.conn = conn
	return nil
}
func (c *TelnetClient) Close() error {
	return c.conn.Close()
}

func (c *TelnetClient) Send() error {
	_, err := io.Copy(c.conn, c.in)
	return err
}

func (c *TelnetClient) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	return err
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
