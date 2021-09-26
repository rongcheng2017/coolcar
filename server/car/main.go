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
	"github.com/streadway/amqp"
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
	amqpConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Fatal("cannot dial amqp", zap.Error(err))
	}
	exchange := "coolcar"

	pub, err := amqpclt.NewPublisher(amqpConn, exchange)
	if err != nil {
		logger.Fatal("cannot create publisher", zap.Error(err))
	}
	carConn, err := grpc.Dial("localhost:8084", grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect car service", zap.Error(err))
	}

	aiConn, err := grpc.Dial("localhost:18001", grpc.WithInsecure())
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
		addr := ":9090"
		logger.Info("Http Server started.", zap.String("Addr", addr))
		logger.Sugar().Fatal(http.ListenAndServe(addr, nil))
	}()
	//Start trip updater.
	tripConn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect rental service", zap.Error(err))
	}
	go trip.RunUpdater(carSub, rentalpb.NewTripServiceClient(tripConn), logger)

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "car",
		Addr:   ":8084",
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
