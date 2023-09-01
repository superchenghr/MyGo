package rpc

import (
	"encoding/gob"
	"fmt"
	"net"
	"testing"
)

func TestRpcClient(t *testing.T) {
	gob.Register(User{})
	addr := "127.0.0.1:8080"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
	}

	cli := NewClient(conn)

	var query func(int) (User, error)

	cli.callRpc("queryUser", &query)

	user, err := query(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("结果....", user)
}
