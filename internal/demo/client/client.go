package client

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	demoservice "github.com/leosunmo/go-grpc-channelz/internal/generated/service"
)

// New creates a new gRPC client
func New(connectionString string) (demoservice.DemoServiceClient, error) {
	return NewWithDialOpts(connectionString,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// NewWithDialOpts creates a new gRPC client with custom []grpc.DialOption
func NewWithDialOpts(connectionString string, dialOpts ...grpc.DialOption) (demoservice.DemoServiceClient, error) {
	conn, err := grpc.Dial(connectionString, dialOpts...)
	if err != nil {
		return nil, errors.Wrapf(err, "error dialing to %s", connectionString)
	}

	client := demoservice.NewDemoServiceClient(conn)
	return client, nil
}
