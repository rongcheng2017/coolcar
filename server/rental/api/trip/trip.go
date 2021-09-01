package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/auth"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// type TripServiceServer interface {
// 	CreateTrip(context.Context, *CreateTripRequest) (*CreateTripResponse, error)
// 	mustEmbedUnimplementedTripServiceServer()
// }

type Service struct {
	Logger *zap.Logger
	rentalpb.UnimplementedTripServiceServer
}

func (s *Service) CreateTrip(c context.Context, request *rentalpb.CreateTripRequest) (*rentalpb.CreateTripResponse, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	s.Logger.Info("create trip", zap.String("start", request.Start), zap.String("account_id", aid.String()))
	return nil, status.Error(codes.Unimplemented, "")
}
