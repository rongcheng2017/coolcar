package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	trippb "coolcar/proto/gen/go"
	trip "coolcar/tripservice"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	trippb.RegisterTripServiceServer(s, &trip.Service{})

	log.Fatal(s.Serve(lis))

}
