package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/auth/wechat"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/shared/server"
	"github.com/namsral/flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", ":8080", "address to listen")
var mongoURI = flag.String("mongo_uri", "mongodb://localhost:27017", "mongo uri to listen")
var privateKeyFile = flag.String("private_key_file", "private.key", "private key file to listen")
var wechatAppID = flag.String("wechat_app_id", "wxcc3d786130252958", "wechat app id")
var wechatAppSecret = flag.String("wechat_app_secret", "ff5309cdd797e3b5e6de46f08a3ab103", "wechat app secret")

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
	}
	pkFile, err := os.Open(*privateKeyFile)
	if err != nil {
		logger.Fatal("cannot open private key", zap.Error(err))
	}
	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key", zap.Error(err))
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
	}

	err = server.RunGRPCServer(&server.GRPCConfig{
		Name:   "auth",
		Logger: logger,
		Addr:   *addr,
		RegisterFunc: func(s *grpc.Server) {
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				OpenIDResolver: &wechat.Service{
					AppID:     *wechatAppID,
					AppSecret: *wechatAppSecret,
				},
				Mongo:          dao.NewMongo(mongoClient.Database("coolcar")),
				Logger:         logger,
				TokenExpire:    2 * time.Hour,
				TokenGenerator: token.NewJWTTokenGen("coolcar/auth", privKey),
			})
		},
	})

	logger.Fatal("cannot start server", zap.Error(err))
}
