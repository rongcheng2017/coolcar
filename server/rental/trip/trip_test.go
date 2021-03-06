package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/client/poi"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"coolcar/shared/server"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestCreateTrip(t *testing.T) {
	//start container
	c := context.Background()
	pm := &profileManager{}
	cm := &carManager{}
	s := newService(c, t, pm, cm)

	req := &rentalpb.CreateTripRequest{
		CarId: "car1",
		Start: &rentalpb.Location{
			Latitude:  32.123,
			Longitude: 114.2525,
		},
	}
	pm.iID = "identity1"
	golden := `{"account_id":%q,"car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":1631095884},"current":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":1631095884},"status":1,"identity_id":"identity1"}`
	nowFunc = func() int64 {
		return 1631095884
	}
	cases := []struct {
		name         string
		accountID    string
		tripID       string
		profileErr   error
		carVerifyErr error
		carUnlockErr error
		want         string
		wantErr      bool
	}{
		{
			name:      "normal_create",
			accountID: "account1",
			tripID:    "61358659c66699ba60b09774",
			want:      fmt.Sprintf(golden, "account1"),
		}, {
			name:       "profile_err",
			accountID:  "account2",
			tripID:     "61358659c66699ba60b09775",
			profileErr: fmt.Errorf("profile"),
			wantErr:    true,
		}, {
			name:         "car_verify_err",
			accountID:    "account3",
			tripID:       "61358659c66699ba60b09776",
			carVerifyErr: fmt.Errorf("verify"),
			wantErr:      true,
		}, {
			name:         "car_unlock_err",
			tripID:       "61358659c66699ba60b09777",
			accountID:    "account4",
			carUnlockErr: fmt.Errorf("unlock"),
			wantErr:      false,
			want:         fmt.Sprintf(golden, "account4"),
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			mgutil.NewObjIDWithValue(id.TripID(cc.tripID))
			pm.err = cc.profileErr
			cm.unlockError = cc.carUnlockErr
			cm.verifyError = cc.carVerifyErr
			c = auth.ContextWithAccountID(context.Background(), id.AccountID(cc.accountID))
			res, err := s.CreateTrip(c, req)
			if cc.wantErr {
				if err == nil {
					t.Errorf("want error; got none")
				} else {
					return
				}
			}
			if err != nil {
				t.Errorf("error creating trip: %v", err)
				return
			}
			if res.Id != cc.tripID {
				t.Errorf("incorrect id; want %q, got %q", cc.tripID, res.Id)
			}
			b, err := json.Marshal(res.Trip)
			if err != nil {
				t.Errorf("cannot marshall response: %v", err)
			}
			got := string(b)
			if cc.want != got {
				t.Errorf("incorrect response: want %s, got %s", cc.want, got)
			}

		})
	}

}

func TestTripLifecycle(t *testing.T) {
	c := auth.ContextWithAccountID(context.Background(), id.AccountID("account_for_lifecycle"))
	s := newService(c, t, &profileManager{}, &carManager{})
	tid := id.TripID("61358659c66699ba60b29777")
	mgutil.NewObjIDWithValue(tid)
	cases := []struct {
		name string
		now  int64
		op   func() (*rentalpb.Trip, error)
		want string
	}{
		{name: "create_trip", now: 10000, op: func() (*rentalpb.Trip, error) {
			e, err := s.CreateTrip(c, &rentalpb.CreateTripRequest{
				CarId: "car1",
				Start: &rentalpb.Location{
					Latitude:  32.123,
					Longitude: 114.2525,
				},
			})
			if err != nil {
				return nil, err
			}
			return e.Trip, nil
		}, want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":10000},"current":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":10000},"status":1}`},
		{name: "update_trip", now: 20000, op: func() (*rentalpb.Trip, error) {
			return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
				Id: tid.String(),
				Current: &rentalpb.Location{
					Latitude:  27.333,
					Longitude: 123.444,
				},
			})
		}, want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":10000},"current":{"location":{"latitude":27.333,"longitude":123.444},"fee_cent":3685,"km_driven":233.60983241807665,"poi_name":"迪士尼","timestamp_sec":20000},"status":1}`},
		{name: "finish_trip", now: 30000, op: func() (*rentalpb.Trip, error) {
			return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
				Id:      tid.String(),
				EndTrip: true,
			})
		}, want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":10000},"end":{"location":{"latitude":27.333,"longitude":123.444},"fee_cent":7914,"km_driven":583.9886135763365,"poi_name":"迪士尼","timestamp_sec":30000},"current":{"location":{"latitude":27.333,"longitude":123.444},"fee_cent":7914,"km_driven":583.9886135763365,"poi_name":"迪士尼","timestamp_sec":30000},"status":2}`},
		{name: "query_trip", now: 40000, op: func() (*rentalpb.Trip, error) {
			return s.GetTrip(c, &rentalpb.GetTripRequest{
				Id: tid.String(),
			})
		}, want: `{"account_id":"account_for_lifecycle","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2525},"poi_name":"天安门","timestamp_sec":10000},"end":{"location":{"latitude":27.333,"longitude":123.444},"fee_cent":7914,"km_driven":583.9886135763365,"poi_name":"迪士尼","timestamp_sec":30000},"current":{"location":{"latitude":27.333,"longitude":123.444},"fee_cent":7914,"km_driven":583.9886135763365,"poi_name":"迪士尼","timestamp_sec":30000},"status":2}`},
	}
	rand.Seed(1345)
	for _, cc := range cases {
		nowFunc = func() int64 {
			return cc.now
		}
		trip, err := cc.op()
		if err != nil {
			t.Errorf("%s: operation failed: %v", cc.name, err)
			continue
		}
		b, err := json.Marshal(trip)
		if err != nil {
			t.Errorf("%s: failed marshlling response: %v", cc.name, err)
		}
		got := string(b)
		if cc.want != got {
			t.Errorf("%s: incorrect response; want: %s, got: %s", cc.name, cc.want, got)
		}

	}
}

func newService(c context.Context, t *testing.T, pm ProfileManager, cm CarManager) *Service {
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongdo db:%v", err)
	}
	logger, err := server.NewZapLogger()
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}
	db := mc.Database("coolcar")
	mongotesting.SetupIndexes(c, db)
	return &Service{
		ProfileManager: pm,
		CarManager:     cm,
		POIManager:     &poi.Manager{},
		Mongo:          dao.NewMongo(db),
		Logger:         logger,
	}
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}

type profileManager struct {
	iID id.IdentityID
	err error
}

func (p *profileManager) Verify(c context.Context, aid id.AccountID) (id.IdentityID, error) {
	return p.iID, p.err
}

type carManager struct {
	verifyError error
	unlockError error
	lockError   error
}

func (cm *carManager) Verify(c context.Context, carID id.CarID, ls *rentalpb.Location) error {
	return cm.verifyError
}

func (cm *carManager) Unlock(ctx context.Context, carID id.CarID, aid id.AccountID, tid id.TripID, avatarURL string) error {
	return cm.unlockError
}
func (cm *carManager) Lock(ctx context.Context, carID id.CarID) error {
	return cm.lockError
}
