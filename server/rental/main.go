package main

import (
	"context"
	blobpb "coolcar/blob/api/gen/v1"
	carpb "coolcar/car/api/gen/v1"
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
	"time"

	"github.com/namsral/flag"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", ":8082", "address to listen")
var blobAddr = flag.String("blob_addr", "localhost:8083", "blob address to listen")
var mongoURI = flag.String("mongo_uri", "mongodb://localhost:27017", "mongo uri ")
var carAddr = flag.String("car_addr", "localhost:8084", "car address")
var authPublicKeyFile = flag.String("auth_public_key_file", `/Users/fengrongcheng/360/golang/coolcar/server/shared/auth/public.key`, "auth public ke file")

func main() {
	flag.Parse()

	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger :%v", err)
	}
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		logger.Fatal("connot connect mongdb", zap.Error(err))
	} else {
		logger.Info("connect mongo db success")
	}
	db := mongoClient.Database("coolcar")
	blobConn, err := grpc.Dial(*blobAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect blob service", zap.Error(err))
	}

	profService := &profile.Service{
		Mongo:             profiledao.NewMongo(db),
		Logger:            logger,
		BlobClient:        blobpb.NewBlobServiceClient(blobConn),
		PhotoGetExpire:    5 * time.Second,
		PhotoUploadExpire: 10 * time.Second,
	}
	carConn, err := grpc.Dial(*carAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect car serivce", zap.Error(err))
	}

	err = server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Logger:            logger,
		Addr:              *addr,
		AuthPublicKeyFile: *authPublicKeyFile,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				CarManager: &car.Manager{
					CarService: carpb.NewCarServiceClient(carConn),
				},
				ProfileManager: &profileClient.Manager{
					Fetcher: profService,
				},
				POIManager: &poi.Manager{},
				Mongo:      tripdao.NewMongo(db),
				Logger:     logger,
			})
			rentalpb.RegisterProfileServiceServer(s, profService)
		},
	})
	logger.Fatal("cannot start trip server", zap.Error(err))

}
