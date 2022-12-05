package server

import (
	api "github.com/izaakdale/service-log/api/v1"
)

type Config struct {
	CommitLog CommitLog
}

var _ api.LogServer = (*grpcServer)(nil)

type grpcServer struct {
	api.UnimplementedLogServer
	*Config
}

func newGrpcServer(cnf *Config) (*grpcServer, error) {
	srv := &grpcServer{
		Config: cnf,
	}
	return srv, nil
}

type CommitLog interface {
	Append(*api.Record) (uint64, error)
	Read(uint64) (*api.Record, error)
}
