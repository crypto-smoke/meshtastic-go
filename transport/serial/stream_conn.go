package serial

import (
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
)

type StreamConn struct {
	conn io.ReadWriteCloser
}

func NewStreamConn(conn io.ReadWriteCloser) *StreamConn {
	return &StreamConn{conn: conn}
}

func (c *StreamConn) Close() (err error) {
	return c.conn.Close()
}

func (c *StreamConn) Read(out proto.Message) error {
	data, err := c.ReadBytes()
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, out)
}

func (c *StreamConn) ReadBytes() ([]byte, error) {
	buf := make([]byte, 4)
	for {
		// Read the first byte, looking for START1.
		_, err := io.ReadFull(c.conn, buf[:1])
		if err != nil {
			return nil, err
		}

		// Check for START1.
		if buf[0] != 0x94 {
			continue
		}

		// Read the second byte, looking for START2.
		_, err = io.ReadFull(c.conn, buf[1:2])
		if err != nil {
			return nil, err
		}

		// Check for START2.
		if buf[1] != 0xc3 {
			continue
		}

		// The next two bytes should be the length of the protobuf message.
		_, err = io.ReadFull(c.conn, buf[2:])
		if err != nil {
			return nil, err
		}

		length := int(binary.BigEndian.Uint16(buf[2:]))
		if length > PACKET_MTU {
			//packet corrupt, start over
			continue
		}
		//fmt.Println("got packet from node with length", length)
		data := make([]byte, length)

		// Read the protobuf data.
		_, err = io.ReadFull(c.conn, data)
		if err != nil {
			return nil, err
		}

		return data, nil
	}
}

func (c *StreamConn) Write(in proto.Message) error {
	protoBytes, err := proto.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshalling proto message: %w", err)
	}
	// TODO: Header
	_, err = c.conn.Write(protoBytes)
}
