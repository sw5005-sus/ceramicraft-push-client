package pushclient

import (
	"fmt"
	"sync"

	"github.com/sw5005-sus/ceramicraft-push-client/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	pushClient pb.NotificationServiceClient
	conn       *grpc.ClientConn
	once       sync.Once
	initErr    error
)

// GetPushClient returns a singleton NotificationServiceClient.
// The connection is established on the first call using the provided address.
// Subsequent calls ignore addr and return the already-initialised client.
func GetPushClient(host string, port int) (pb.NotificationServiceClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	once.Do(func() {
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(1024*1024),
				grpc.MaxCallSendMsgSize(1024*1024),
			),
		}
		var err error
		conn, err = grpc.NewClient(addr, opts...)
		if err != nil {
			initErr = fmt.Errorf("pushclient: failed to connect to %s: %w", addr, err)
			return
		}
		pushClient = pb.NewNotificationServiceClient(conn)
	})
	return pushClient, initErr
}

// Close closes the underlying gRPC connection. It should be called when the
// client is no longer needed, typically during application shutdown.
func Close() error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}
