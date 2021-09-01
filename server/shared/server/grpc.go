package server

import (
	"coolcar/shared/auth"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCConfig struct {
	Name              string
	Logger            *zap.Logger
	Addr              string
	AuthPublicKeyFile string
	RegisterFunc      func(*grpc.Server)
}

func RunGRPCServer(c *GRPCConfig) error {
	nameField := zap.String("name", c.Name)
	var logger = c.Logger

	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		logger.Fatal("cannot listen", nameField, zap.Error(err))
	}

	var opts []grpc.ServerOption
	if c.AuthPublicKeyFile != "" {
		in, err := auth.Inteceptor(c.AuthPublicKeyFile)
		if err != nil {
			logger.Fatal("cannot create auth interceptor", zap.Error(err))
		}
		opts = append(opts, grpc.UnaryInterceptor(in))
	}

	s := grpc.NewServer(opts...)

	c.RegisterFunc(s)

	logger.Info("server started", nameField, zap.String("addr", c.Addr))
	return s.Serve(lis)
}
