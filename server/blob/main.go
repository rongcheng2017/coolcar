package main

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	"coolcar/blob/blob"
	"coolcar/blob/blob/cos"
	"coolcar/blob/blob/dao"
	"coolcar/shared/server"
	"log"

	"github.com/namsral/flag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", ":8083", "address to listen")
var mongoURI = flag.String("mongo_uri", "mongodb://localhost:27017", "mongo uri to listen")
var cosAddr = flag.String("cos_addr", "https://coolcar-1307431695.cos.ap-beijing.myqcloud.com", "cos uri")
var cosSecID = flag.String("cos_sec_id", "", "cos sercet id")
var cosSecKey = flag.String("cos_sec_key", "", "cos sercet key")

func main() {
	flag.Parse()

	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}
	db := mongoClient.Database("coolcar")
	st, err := cos.NewService(*cosAddr, *cosSecID, *cosSecKey)
	if err != nil {
		logger.Fatal("cannot create cos service", zap.Error(err))
	}
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "blob",
		Addr:   *addr,
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
