package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	carpb "coolcar/car/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/auth"
	"coolcar/shared/server"
	"github.com/namsral/flag"
	"log"
	"net/http"
	"net/textproto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

var addr = flag.String("addr", ":8080", "address to listen")
var authAddr = flag.String("auth_addr", ":loalhost:8081", "address to listen")
var tripAddr = flag.String("trip_addr", "loalhost:8082", "address to listen")
var profileAddr = flag.String("profile_addr", "loalhost:8082", "address to listen")
var carAddr = flag.String("car_addr", "loalhost:8084", "address to listen")

func main() {
	flag.Parse()
	
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
			addr:         *authAddr,
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "trip",
			addr:         *tripAddr,
			registerFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		}, {
			name:         "profile",
			addr:         *profileAddr,
			registerFunc: rentalpb.RegisterProfileServiceHandlerFromEndpoint,
		},
		{
			name:         "car",
			addr:         *carAddr,
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
	logger.Sugar().Infof("grpc gateway started at %s", *addr)
	//listen and serve
	logger.Sugar().Fatal(http.ListenAndServe(*addr, mux))

}
