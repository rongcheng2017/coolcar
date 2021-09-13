package main

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	"coolcar/blob/blob"
	"coolcar/blob/blob/cos"
	"coolcar/blob/blob/dao"
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
		log.Fatalf("cannot create logger: %v", err)
	}
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&appname=mongodb-vscode%200.6.10&directConnection=true&ssl=false"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}
	db := mongoClient.Database("coolcar")
	st, err := cos.NewService("https://coolcar-1307431695.cos.ap-beijing.myqcloud.com", "AKIDVpseIflbCYeT2KL8gxxg8KFtHPGq9CyB", "0kFSO3SaQ1OUKtBE1bvzmJZaJlxsXv1q")
	if err != nil {
		logger.Fatal("cannot create cos service", zap.Error(err))
	}
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "blob",
		Addr:   ":8083",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			blobpb.RegisterBlobServiceServer(s, &blob.Service{
				Storage: st,
				Mongo:   dao.NewMongo(db),
				Logger:  logger,
			})
		},
	}))

}
