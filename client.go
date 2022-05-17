package channelz

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	channelzgrpc "google.golang.org/grpc/channelz/grpc_channelz_v1"
	"google.golang.org/grpc/credentials/insecure"
)

func (h *grpcChannelzHandler) connect() (channelzgrpc.ChannelzClient, error) {
	if h.client != nil {
		// Already connected
		return h.client, nil
	}

	host := getHostFromBindAddress(h.bindAddress)
	h.mu.Lock()
	defer h.mu.Unlock()
	client, err := newChannelzClient(host)
	if err != nil {
		return nil, err
	}
	h.client = client
	return h.client, nil
}

func newChannelzClient(dialString string) (channelzgrpc.ChannelzClient, error) {
	conn, err := grpc.Dial(dialString,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error dialing to %s", dialString)
	}
	client := channelzgrpc.NewChannelzClient(conn)
	return client, nil
}

func getHostFromBindAddress(bindAddress string) string {
	if strings.HasPrefix(bindAddress, ":") {
		return fmt.Sprintf("localhost%s", bindAddress)
	}
	return bindAddress
}
