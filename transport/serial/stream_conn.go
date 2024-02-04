package serial

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
)

// StreamConn implements the meshtastic client API stream protocol.
// This protocol is used to send and receive protobuf messages over a serial or TCP connection.
// See https://meshtastic.org/docs/development/device/client-api#streaming-version for additional information.
type StreamConn struct {
	conn io.ReadWriteCloser
}

// NewStreamConn creates a new StreamConn with the provided io.ReadWriteCloser.
// Once an io.ReadWriteCloser is provided, the StreamConn should be used read, write and close operations.
func NewStreamConn(conn io.ReadWriteCloser) *StreamConn {
	return &StreamConn{conn: conn}
}

// Close closes the connection.
func (c *StreamConn) Close() (err error) {
	return c.conn.Close()
}

// Read reads a protobuf message from the connection.
func (c *StreamConn) Read(out proto.Message) error {
	data, err := c.ReadBytes()
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, out)
}

// ReadBytes reads a byte slice from the connection.
// Prefer using Read if you have a protobuf message.
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

// writeStreamHeader writes the stream protocol header to the provided writer.
// See https://meshtastic.org/docs/development/device/client-api#streaming-version
func writeStreamHeader(w io.Writer, dataLen uint16) error {
	header := bytes.NewBuffer(nil)
	// First we write START1, START2
	header.WriteByte(START1)
	header.WriteByte(START2)
	// Next we write the length of the protobuf message as a big-endian uint16
	err := binary.Write(header, binary.BigEndian, dataLen)
	if err != nil {
		return fmt.Errorf("writing length to buffer: %w", err)
	}

	_, err = w.Write(header.Bytes())
	return err
}

// Write writes a protobuf message to the connection.
func (c *StreamConn) Write(in proto.Message) error {
	protoBytes, err := proto.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshalling proto message: %w", err)
	}

	if err := c.WriteBytes(protoBytes); err != nil {
		return err
	}

	return nil
}

// WriteBytes writes a byte slice to the connection.
// Prefer using Write if you have a protobuf message.
func (c *StreamConn) WriteBytes(data []byte) error {
	// TODO: Lock this function to prevent concurrent writes
	if err := writeStreamHeader(c.conn, uint16(len(data))); err != nil {
		return fmt.Errorf("writing stream header: %w", err)
	}

	if _, err := c.conn.Write(data); err != nil {
		return fmt.Errorf("writing proto message: %w", err)
	}
	return nil
}
