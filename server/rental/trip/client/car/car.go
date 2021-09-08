package car

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
)

type Manager struct {
}


func (cm *Manager) Verify(ctx context.Context, carID id.CarID, loc *rentalpb.Location) error {
	return nil
}

func (cm *Manager) Unlock(ctx context.Context, carID id.CarID) error {
	return nil
}
