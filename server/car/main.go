package main

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	"coolcar/car/car"
	"coolcar/car/dao"
	amqpclt "coolcar/car/mq/amqpclt"
	"coolcar/car/sim"
	"coolcar/car/sim/pos"
	"coolcar/car/trip"
	"coolcar/car/ws"
	rentalpb "coolcar/rental/api/gen/v1"
	coolenvpb "coolcar/shared/coolenv"
	"coolcar/shared/server"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/namsral/flag"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", ":8084", "address to listen")
var wsAddr = flag.String("ws_addr", ":9090", "weksocker address to listen")
var mongoURI = flag.String("mongo_uri", "mongodb://localhost:27017", "mongo uri ")
var amqpURL = flag.String("amqp_url", "amqp://guest:guest@localhost:5672/", "amqp url ")
var carAddr = flag.String("car_addr", "localhost:8084", "car address")
var tripAddr = flag.String("trip_addr", "localhost:8082", "trip address")
var aiAddr = flag.String("ai_addr", "localhost:18001", "ai address")

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
	amqpConn, err := amqp.Dial(*amqpURL)
	if err != nil {
		logger.Fatal("cannot dial amqp", zap.Error(err))
	}
	exchange := "coolcar"

	pub, err := amqpclt.NewPublisher(amqpConn, exchange)
	if err != nil {
		logger.Fatal("cannot create publisher", zap.Error(err))
	}
	carConn, err := grpc.Dial(*carAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect car service", zap.Error(err))
	}

	aiConn, err := grpc.Dial(*aiAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect ai service", zap.Error(err))
	}

	carSub, err := amqpclt.NewSubscriber(amqpConn, exchange, logger)
	if err != nil {
		logger.Fatal("cannot create car subscriber", zap.Error(err))
	}
	poSub, err := amqpclt.NewSubscriber(amqpConn, "pos_sim", logger)
	if err != nil {
		logger.Fatal("cannot create pos subscriber", zap.Error(err))
	}

	simController := sim.Controller{
		CarService:    carpb.NewCarServiceClient(carConn),
		CarSubscriber: carSub,
		AIService:     coolenvpb.NewAIServiceClient(aiConn),
		PosSubScriber: &pos.Subscriber{
			Sub:    poSub,
			Logger: logger,
		},
		Logger: logger,
	}

	go simController.RunSimulations(context.Background())

	//Start websocket handler.
	u := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	http.HandleFunc("/ws", ws.Handler(u, carSub, logger))

	go func() {
		addr := *wsAddr
		logger.Info("Http Server started.", zap.String("Addr", addr))
		logger.Sugar().Fatal(http.ListenAndServe(addr, nil))
	}()
	//Start trip updater.
	tripConn, err := grpc.Dial(*tripAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect rental service", zap.Error(err))
	}
	go trip.RunUpdater(carSub, rentalpb.NewTripServiceClient(tripConn), logger)

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "car",
		Addr:   *addr,
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			carpb.RegisterCarServiceServer(s, &car.Service{
				Mongo:     dao.NewMongo(db),
				Logger:    logger,
				Publisher: pub,
			})
		},
	}))

}
