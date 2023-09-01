package rpc

import (
	"encoding/gob"
	"fmt"
	"sync"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func queryUser(uid int) (User, error) {
	user := make(map[int]User)
	user[0] = User{"zs", 20}
	user[1] = User{"ls", 21}
	if u, ok := user[uid]; ok {
		return u, nil
	}
	return User{}, fmt.Errorf("user not found")
}

func TestRpc(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	gob.Register(User{})
	addr := "127.0.0.1:8080"
	srv := NewSever(addr)
	srv.Register("queryUser", queryUser)
	go srv.Run()
	fmt.Println("服务端启动成功。。。。。。")
	wg.Wait()
}
