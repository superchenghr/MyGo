package rpc

import (
	"fmt"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) callRpc(rpcName string, fPtr interface{}) {
	fn := reflect.ValueOf(fPtr).Elem()
	f := func(args []reflect.Value) []reflect.Value {
		inArgs := make([]interface{}, 0, len(args))
		for _, arg := range args {
			inArgs = append(inArgs, arg.Interface())
		}
		cliSession := newSession(c.conn)
		req := RPCData{Name: rpcName, Args: inArgs}
		bytes, err := encode(req)
		if err != nil {
			panic(err)
		}
		err = cliSession.Write(bytes)
		if err != nil {
			panic(err)
		}
		resp, err := cliSession.read()
		respData, err := decode(resp)
		if err != nil {
			panic(err)
		}
		values := make([]reflect.Value, 0, len(respData.Args))
		for i, arg := range respData.Args {
			if arg == nil {
				values = append(values, reflect.Zero(fn.Type().Out(i)))
				continue
			}
			values = append(values, reflect.ValueOf(arg))
		}
		return values
	}
	value := reflect.MakeFunc(fn.Type(), f)
	fmt.Println(fn.Type())
	fn.Set(value)
}
