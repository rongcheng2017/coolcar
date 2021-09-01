package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/api/trip"
	"coolcar/shared/server"
	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger :%v", err)
	}

	err = server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Logger:            logger,
		Addr:              ":8082",
		AuthPublicKeyFile: "/Users/fengrongcheng/360/golang/coolcar/server/shared/auth/public.key",
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				Logger: logger,
			})
		},
	})
	logger.Fatal("cannot start server", zap.Error(err))
}
