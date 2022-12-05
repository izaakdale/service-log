package server

import (
	api "github.com/izaakdale/service-log/api/v1"
	"google.golang.org/grpc"
)

type Config struct {
	CommitLog CommitLog
}

var _ api.LogServer = (*grpcServer)(nil)

func NewGRPCServer(cfg *Config) (*grpc.Server, error) {
	gsrv := grpc.NewServer()
	srv, err := newGrpcServer(cfg)
	if err != nil {
		return nil, err
	}
	api.RegisterLogServer(gsrv, srv)
	return gsrv, nil
}

type grpcServer struct {
	api.UnimplementedLogServer
	*Config
}

func newGrpcServer(cfg *Config) (*grpcServer, error) {
	srv := &grpcServer{
		Config: cfg,
	}
	return srv, nil
}

type CommitLog interface {
	Append(*api.Record) (uint64, error)
	Read(uint64) (*api.Record, error)
}
