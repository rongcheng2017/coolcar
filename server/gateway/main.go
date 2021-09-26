package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	carpb "coolcar/car/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/auth"
	"coolcar/shared/server"
	"log"
	"net/http"
	"net/textproto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger :%v", err)
	}

	//上下文
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()
	//ServeMux
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseEnumNumbers: true,
				UseProtoNames:  true,
			},
		},
	), runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
		if s == textproto.CanonicalMIMEHeaderKey(runtime.MetadataHeaderPrefix+auth.ImpersonateAccountHeader) {
			return "", false
		}
		return runtime.DefaultHeaderMatcher(s)
	}))

	serverConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         "localhost:8081",
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "trip",
			addr:         "localhost:8082",
			registerFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		}, {
			name:         "profile",
			addr:         "localhost:8082",
			registerFunc: rentalpb.RegisterProfileServiceHandlerFromEndpoint,
		},
		{
			name:         "car",
			addr:         "localhost:8084",
			registerFunc: carpb.RegisterCarServiceHandlerFromEndpoint,
		},
	}
	// register servers
	for _, sc := range serverConfig {
		err := sc.registerFunc(c, mux, sc.addr, []grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			logger.Sugar().Fatalf("cannot register %s  service %v", sc.name, err)
		}
	}
	addr := ":8080"
	logger.Sugar().Infof("grpc gateway started at %s", addr)
	//listen and serve
	logger.Sugar().Fatal(http.ListenAndServe(addr, mux))

}
