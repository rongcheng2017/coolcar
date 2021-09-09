package main

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/profile"
	profiledao "coolcar/rental/profile/dao"
	"coolcar/rental/trip"
	"coolcar/rental/trip/client/car"
	"coolcar/rental/trip/client/poi"
	profileClient "coolcar/rental/trip/client/profile"
	tripdao "coolcar/rental/trip/dao"
	"coolcar/shared/server"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger :%v", err)
	}
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&appname=mongodb-vscode%200.6.10&directConnection=true&ssl=false"))
	if err != nil {
		logger.Fatal("connot connect mongdb", zap.Error(err))
	} else {
		logger.Info("connect mongo db success")
	}
	db := mongoClient.Database("coolcar")
	err = server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Logger:            logger,
		Addr:              ":8082",
		AuthPublicKeyFile: "/Users/fengrongcheng/360/golang/coolcar/server/shared/auth/public.key",
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				CarManager:     &car.Manager{},
				ProfileManager: &profileClient.Manager{},
				POIManager:     &poi.Manager{},
				Mongo:          tripdao.NewMongo(db),
				Logger:         logger,
			})
			rentalpb.RegisterProfileServiceServer(s, &profile.Service{
				Mongo:  profiledao.NewMongo(db),
				Logger: logger,
			})
		},
	})
	logger.Fatal("cannot start trip server", zap.Error(err))

}
