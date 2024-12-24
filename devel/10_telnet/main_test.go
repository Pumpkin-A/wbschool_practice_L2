package main

import (
	"net"
	"sync"
	"testing"
	"time"
)

func TestTelnetClient(t *testing.T) {
	t.Run("defaule", func(t *testing.T) {
		l, _ := net.Listen("tcp", "127.0.0.1:")
		defer l.Close()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			client := NewTelnetClient(l.Addr().String(), time.Second*1)
			if err := client.initConnection(); err != nil {
				t.Error(err.Error())
			}
			defer func() {
				if err := client.closeConnection(); err != nil {
					t.Error(err.Error())
				}
			}()

			if err := client.sendMsg("hello\n"); err != nil {
				t.Errorf("error! send msg should not provied error")
			}
		}()

		go func() {
			defer wg.Done()

			conn, _ := l.Accept()
			defer conn.Close()

			request := make([]byte, 128)
			readed, _ := conn.Read(request)
			if "hello\n" != string(request)[:readed] {
				t.Errorf("Error! Msgs not equal")
			}
		}()

		wg.Wait()
	})

	t.Run("should return error when connect to wrong host", func(t *testing.T) {
		client := NewTelnetClient("localhost:909090", time.Second*1)
		if err := client.initConnection(); err == nil {
			t.Errorf("wrong addres should error appear")
		}
	})
}
