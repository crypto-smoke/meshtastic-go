package serial

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestStreamConn(t *testing.T) {
	deviceNetConn, clientNetConn := net.Pipe()
	_ = NewStreamConn(deviceNetConn)
	_ = NewStreamConn(clientNetConn)
}

func Test_writeStreamHeader(t *testing.T) {
	out := bytes.NewBuffer(nil)
	err := writeStreamHeader(out, 257)
	require.NoError(t, err)
	require.Equal(t, []byte{START1, START2, 0x01, 0x01}, out.Bytes())
}
