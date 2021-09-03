package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"

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

func (s *Service) CreateTrip(c context.Context, request *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	return nil, status.Error(codes.Unauthenticated, "")
}
func (s *Service) GetTrip(c context.Context, request *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	return nil, status.Error(codes.Unauthenticated, "")
}
func (s *Service) GetTrips(c context.Context, request *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
	return nil, status.Error(codes.Unauthenticated, "")
}
func (s *Service) UpdateTrip(c context.Context, request *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {
	return nil, status.Error(codes.Unauthenticated, "")
}
