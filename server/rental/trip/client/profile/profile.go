package profile

import (
	"context"
	"coolcar/shared/id"
)

type Manager struct {
}

func (pm *Manager) Verify(ctx context.Context, aid id.AccountID) (id.IdentityID, error) {
	return id.IdentityID("identity1"), nil
}
