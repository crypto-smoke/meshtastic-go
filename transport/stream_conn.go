package transport

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"sync"
	"time"
)

const (
	// WaitAfterWake is the amount of time to wait after sending the wake message before sending the first message.
	WaitAfterWake = 100 * time.Millisecond
	// Start1 is a magic byte used in the meshtastic stream protocol.
	// Start1 is sent at the beginning of a message to indicate the start of a new message.
	Start1 = 0x94
	// Start2 is a magic byte used in the meshtastic stream protocol.
	// It is sent after Start1 to indicate the start of a new message.
	Start2 = 0xc3
	// PacketMTU is the maximum size of the protobuf message which can be sent within the header.
	PacketMTU = 512
)

// StreamConn implements the meshtastic client API stream protocol.
// This protocol is used to send and receive protobuf messages over a serial or TCP connection.
// See https://meshtastic.org/docs/development/device/client-api#streaming-version for additional information.
type StreamConn struct {
	conn io.ReadWriteCloser
	// DebugWriter is an optional writer that is used when a non-protobuf message is sent over the connection.
	DebugWriter io.Writer

	readMu  sync.Mutex
	writeMu sync.Mutex
}

// NewClientStreamConn creates a new StreamConn with the provided io.ReadWriteCloser.
// Once an io.ReadWriteCloser is provided, the StreamConn should be used read, write and close operations.
func NewClientStreamConn(conn io.ReadWriteCloser) (*StreamConn, error) {
	sConn := &StreamConn{conn: conn}
	if err := sConn.writeWake(); err != nil {
		return nil, fmt.Errorf("sending wake message: %w", err)
	}
	return sConn, nil
}

// NewRadioStreamConn creates a new StreamConn with the provided io.ReadWriteCloser.
// Once an io.ReadWriteCloser is provided, the StreamConn should be used read, write and close operations.
func NewRadioStreamConn(conn io.ReadWriteCloser) *StreamConn {
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

// ReadBytes reads a byte message from the connection.
// Prefer using Read if you have a protobuf message.
func (c *StreamConn) ReadBytes() ([]byte, error) {
	c.readMu.Lock()
	defer c.readMu.Unlock()
	buf := make([]byte, 4)
	for {
		// Read the first byte, looking for Start1.
		_, err := io.ReadFull(c.conn, buf[:1])
		if err != nil {
			return nil, err
		}

		// Check for Start1.
		if buf[0] != Start1 {
			if c.DebugWriter != nil {
				c.DebugWriter.Write(buf[0:1])
			}
			continue
		}

		// Read the second byte, looking for Start2.
		_, err = io.ReadFull(c.conn, buf[1:2])
		if err != nil {
			return nil, err
		}

		// Check for Start2.
		if buf[1] != Start2 {
			continue
		}

		// The next two bytes should be the length of the protobuf message.
		_, err = io.ReadFull(c.conn, buf[2:])
		if err != nil {
			return nil, err
		}

		length := int(binary.BigEndian.Uint16(buf[2:]))
		if length > PacketMTU {
			//packet corrupt, start over
			continue
		}
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
	// First we write Start1, Start2
	header.WriteByte(Start1)
	header.WriteByte(Start2)
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
	if len(data) > PacketMTU {
		return fmt.Errorf("data length exceeds MTU: %d > %d", len(data), PacketMTU)
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if err := writeStreamHeader(c.conn, uint16(len(data))); err != nil {
		return fmt.Errorf("writing stream header: %w", err)
	}

	if _, err := c.conn.Write(data); err != nil {
		return fmt.Errorf("writing proto message: %w", err)
	}
	return nil
}

// writeWake writes a wake message to the radio.
// This should only be called on the client side.
//
// TODO: Rather than just sending this on start, do we need to also send this after a long period of inactivity?
func (c *StreamConn) writeWake() error {
	// Send 32 bytes of Start2 to wake the radio if sleeping.
	_, err := c.conn.Write(
		bytes.Repeat([]byte{Start2}, 32),
	)
	if err != nil {
		return err
	}
	time.Sleep(WaitAfterWake)
	return nil
}
