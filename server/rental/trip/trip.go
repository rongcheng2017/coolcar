package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Logger         *zap.Logger
	Mongo          *dao.Mongo
	ProfileManager ProfileManager
	CarManager     CarManager
	POIManager     POIManager
	rentalpb.UnimplementedTripServiceServer
}

// ProfileManager defines the ACL(Anti Corruption Layer)
type ProfileManager interface {
	Verify(context.Context, id.AccountID) (id.IdentityID, error)
}

//CarManager defines the ACL for car management
type CarManager interface {
	//需要知道车的位置以及人的位置
	Verify(context.Context, id.CarID, *rentalpb.Location) error
	Unlock(context.Context, id.CarID) error
}

//POIManager resolves Ponit of Interest
type POIManager interface {
	Resolve(context.Context, *rentalpb.Location) (string, error)
}

func (s *Service) CreateTrip(c context.Context, request *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	s.Logger.Info("received create trip cmd ")

	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}else{
		s.Logger.Sugar().Infof("trip : get aid %s",aid)
	}
	//验证驾驶者身份，预防验证身份后，用户修改了身份信息。iID可用来追溯
	iID, err := s.ProfileManager.Verify(c, aid)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}
	//检测车辆状态
	carID := id.CarID(request.CarId)
	err = s.CarManager.Verify(c, carID, request.Start)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	poi, err := s.POIManager.Resolve(c, request.Start)
	if err != nil {
		s.Logger.Info("cannot resolve poi", zap.Stringer("location", request.Start), zap.Error(err))
	}
	//创建行程：写入数据库，开始计费  --> 先创建行程再开锁有补救措施，先开锁再创建行程(创建行程失败)则完犊子了。
	ls := &rentalpb.LocationStatus{
		Location: request.Start,
		PoiName:  poi,
	}
	tr, err := s.Mongo.CreateTrip(c, &rentalpb.Trip{
		AccountId:  aid.String(),
		CarId:      carID.String(),
		IdentityId: iID.String(),
		Status:     rentalpb.TripStatus_IN_PROGRESS,
		Start:      ls,
		Current:    ls,
	})
	if err != nil {
		s.Logger.Warn("cannot create trip", zap.Error(err))
		return nil, status.Error(codes.AlreadyExists, "")
	}

	//车辆开锁
	go func() {
		err := s.CarManager.Unlock(context.Background(), carID)
		if err != nil {
			s.Logger.Error("cannot unlock car ", zap.Error(err))
		}
	}()

	return &rentalpb.TripEntity{
		Id:   tr.ID.Hex(),
		Trip: tr.Trip,
	}, nil
}





func (s *Service) GetTrip(c context.Context, request *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	s.Logger.Info("received GetTrip  cmd ")
	return nil, status.Error(codes.Unauthenticated, "")
}
func (s *Service) GetTrips(c context.Context, request *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
	s.Logger.Info("received GetTrips  cmd ")
	return nil, status.Error(codes.Unauthenticated, "")
}
func (s *Service) UpdateTrip(c context.Context, request *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {
	s.Logger.Info("received UpdateTrip  cmd ")
	return nil, status.Error(codes.Unauthenticated, "")
}
