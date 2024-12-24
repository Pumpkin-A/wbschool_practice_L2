package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

// TelnetClient do basic support of sending and receiving via telnet
type TelnetClient struct {
	address    string
	timeout    time.Duration
	conn       net.Conn
	connReader *bufio.Reader
}

// NewTelnetClient creates instance of TelnetClient
func NewTelnetClient(address string, timeout time.Duration) *TelnetClient {
	res := &TelnetClient{
		address: address,
		timeout: timeout,
	}
	return res
}

func (t *TelnetClient) initConnection() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("%s : %w", "Conncetion error", err)
	}
	t.conn = conn
	t.connReader = bufio.NewReader(t.conn)
	fmt.Printf("Connected success to  %v\n", t.address)
	return nil
}

func (t *TelnetClient) closeConnection() error {
	return t.conn.Close()
}

func (t *TelnetClient) receieveMsg() error {
	line, err := t.connReader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			err = fmt.Errorf("error: closed by peer")
		}
		return err
	}
	if _, err := fmt.Print(line); err != nil {
		return err
	}
	return nil
}

func (t *TelnetClient) sendMsg(msg string) error {
	_, err := t.conn.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("error: closed by peer")
	}
	return nil
}
