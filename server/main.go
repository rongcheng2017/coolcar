package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	trippb "coolcar/proto/gen/go"
	trip "coolcar/tripservice"
)

func main() {
	log.SetFlags(log.Lshortfile)
	go startGRPCGateWay()

	//grpc server
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	trippb.RegisterTripServiceServer(s, &trip.Service{})

	log.Fatal(s.Serve(lis))

}

func startGRPCGateWay() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	//close
	defer cancel()
	//grpc<-->json的过程中
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				//gateway转换过程中将枚举换成对应number
				UseEnumNumbers: true,
				//不适用驼峰命名
				UseProtoNames: true,
			},
		},
	))
	err := trippb.RegisterTripServiceHandlerFromEndpoint(
		c,
		mux,
		"localhost:8081",
		[]grpc.DialOption{grpc.WithInsecure()},
	)
	if err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}

	er := http.ListenAndServe(":8080", mux)
	if er != nil {
		log.Fatalf("cannot http ListenAndServe: %v", er)
	}

}
