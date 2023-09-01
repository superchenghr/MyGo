package rpc

import (
	"encoding/binary"
	"io"
	"net"
)

type Session struct {
	conn net.Conn
}

func newSession(conn net.Conn) *Session {
	return &Session{conn: conn}
}

func (s *Session) Write(data []byte) error {
	buf := make([]byte, 4+len(data))
	// tcp前四位存Type、Len、Value，Type
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)
	_, err := s.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) read() ([]byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(s.conn, header)
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(s.conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
