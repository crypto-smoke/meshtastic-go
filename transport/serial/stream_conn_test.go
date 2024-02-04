package serial

import (
	pb "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"bytes"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"net"
	"testing"
)

func TestStreamConn(t *testing.T) {
	radioNetConn, clientNetConn := net.Pipe()

	sent := &pb.ToRadio{
		PayloadVariant: &pb.ToRadio_WantConfigId{
			WantConfigId: 123,
		},
	}
	received := &pb.ToRadio{}

	eg := errgroup.Group{}
	eg.Go(func() error {
		client, err := NewClientStreamConn(clientNetConn)
		require.NoError(t, err)
		return client.Write(sent)
	})
	eg.Go(func() error {
		radio := NewRadioStreamConn(radioNetConn)
		return radio.Read(received)
	})

	require.NoError(t, eg.Wait())
	require.True(t, proto.Equal(sent, received))
}

func Test_writeStreamHeader(t *testing.T) {
	out := bytes.NewBuffer(nil)
	err := writeStreamHeader(out, 257)
	require.NoError(t, err)
	require.Equal(t, []byte{Start1, Start2, 0x01, 0x01}, out.Bytes())
}
