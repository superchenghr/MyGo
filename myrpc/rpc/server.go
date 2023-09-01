package rpc

import (
	"fmt"
	"net"
	"reflect"
)

type Server struct {
	addr  string
	funcs map[string]reflect.Value
}

func NewSever(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

func (s *Server) Register(rpcName string, f interface{}) {
	if _, ok := s.funcs[rpcName]; ok {
		return
	}
	fVal := reflect.ValueOf(f)
	s.funcs[rpcName] = fVal
}

func (s *Server) Run() {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Println("listen error")
		return
	}
	for {
		fmt.Println("等待接收。。。。")
		conn, err := listen.Accept()
		fmt.Println("接收到请求，开始处理。。")
		if err != nil {
			fmt.Println("accept error")
			return
		}
		session := newSession(conn)
		b, err := session.read()
		if err != nil {
			fmt.Println("read error")
			return
		}
		data, err := decode(b)
		if err != nil {
			fmt.Println("decode error")
			return
		}
		f, ok := s.funcs[data.Name]
		if !ok {
			fmt.Println("func not exist")
			return
		}
		values := make([]reflect.Value, 0, len(data.Args))
		for _, arg := range data.Args {
			values = append(values, reflect.ValueOf(arg))
		}
		out := f.Call(values)

		outArgs := make([]interface{}, 0, len(out))

		for _, o := range out {
			outArgs = append(outArgs, o.Interface())
		}
		resp := RPCData{data.Name, outArgs}
		respBytes, err := encode(resp)
		if err != nil {
			fmt.Println("encode error")
			return
		}
		err = session.Write(respBytes)
		if err != nil {
			fmt.Println("write error")
			return
		}
	}
}
