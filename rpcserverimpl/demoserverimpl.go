package rpcserverimpl

import (
	"context"
	"grpclb/rpcfile"
)

// DemoServiceServer is the server API for DemoService service.
type DemoServiceServer struct {
	ServerAddress string
}

func (s *DemoServiceServer) DemoHandler(_ context.Context, r *demo.DemoRequest) (*demo.DemoResponse, error) {
	return &demo.DemoResponse{Name: r.Name + s.ServerAddress}, nil
}
