package car

import (
	"context"
	carpb "coolcar/car/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	"fmt"
)

type Manager struct {
	CarService carpb.CarServiceClient
}

func (m *Manager) Verify(ctx context.Context, carID id.CarID, loc *rentalpb.Location) error {
	car, err := m.CarService.GetCar(ctx, &carpb.GetCarRequest{
		Id: carID.String(),
	})
	if err != nil {
		return fmt.Errorf("cannot get car: %v", err)
	}
	if car.Status != carpb.CarStatus_LOCKED {
		return fmt.Errorf("cannot unlock; car status is %v", car.Status)
	}
	return nil
}

func (m *Manager) Unlock(ctx context.Context, carID id.CarID, aid id.AccountID, tid id.TripID, avatarURL string) error {
	_, err := m.CarService.UnlockCar(ctx, &carpb.UnlockCarRequest{
		Id: string(carID),
		Driver: &carpb.Driver{
			Id:        string(aid),
			AvatarUrl: avatarURL,
		},
		TripId: string(tid),
	})
	if err != nil {
		return fmt.Errorf("cannot unlock: %v", err)
	}
	return nil
}

func (m *Manager) Lock(c context.Context, carID id.CarID) error {
	_, err := m.CarService.LockCar(c, &carpb.LockCarRequest{
		Id: carID.String(),
	})
	if err != nil {
		return fmt.Errorf("cannot lock: %v", err)
	}
	return nil
}
