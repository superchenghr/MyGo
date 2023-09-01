package rpc

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestSession_Write(t *testing.T) {
	addr := "127.0.0.1:8000"
	data := "hello"
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		conn, _ := lis.Accept()
		s := Session{conn: conn}
		err = s.Write([]byte(data))
		if err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		s := Session{conn: conn}
		data, err := s.read()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(data))
	}()
	wg.Wait()
}
