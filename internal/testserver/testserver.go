package testserver

import (
	"context"
)

var _ TestServerServer = (*EchoTestServer)(nil)

type EchoTestServer struct{}

func (f EchoTestServer) SendMessage(ctx context.Context, request *MessageRequest) (*MessageReply, error) {
	return &MessageReply{Response: request.Request}, nil
}
