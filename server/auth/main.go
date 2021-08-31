package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/auth/wechat"
	"coolcar/auth/dao"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := newZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger :%v", err)
	}

	lis, er := net.Listen("tcp", ":8081")
	if er != nil {
		logger.Fatal("cannot listen", zap.Error(er))
	}

	c:= context.Background()
	mongoClient,err := mongo.Connect(c,options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&appname=mongodb-vscode%200.6.10&directConnection=true&ssl=false"))
	if err != nil {
		logger.Fatal("connot connect mongdb",zap.Error(err))
	}

	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		OpenIDResolver: &wechat.Service{
			AppID: "wxcc3d786130252958",
			AppSecret: "ff5309cdd797e3b5e6de46f08a3ab103",
		},
		Mongo: dao.NewMongo(mongoClient.Database("coolcar")),
		Logger: logger,
	})

	err =s.Serve(lis)
	logger.Fatal("cannot server",zap.Error(err))
}

func newZapLogger() (*zap.Logger,error){
	cfg:= zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey=""
	return cfg.Build()
}