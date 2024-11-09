package manor

import (
	"context"

	"github.com/light-speak/lighthouse/manor/kitex_gen/manor/rpc"
)

type ManorImpl struct{}

func (s *ManorImpl) Register(ctx context.Context, req *rpc.RegisterRequest) (*rpc.RegisterResponse, error) {
	err := RegisterService(req.ServiceName, req.ServiceAddr, req.Store)
	if err != nil {
		return &rpc.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &rpc.RegisterResponse{
		Success: true,
		Message: "Service registered successfully",
	}, nil
}

func (s *ManorImpl) Ping(ctx context.Context, req *rpc.PingRequest) (*rpc.PingResponse, error) {
	return &rpc.PingResponse{
		Message: "pong",
	}, nil
}
