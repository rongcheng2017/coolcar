package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"os"

	mongotesting "coolcar/shared/mongo/testing"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCreateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot create mongo db :%v", err)
	}

	db := mc.Database("coolcar")
	err = mongotesting.SetupIndexes(c, db)
	if err != nil {
		t.Fatalf("cannot setup indexes: %v", err)
	}

	m := NewMongo(db)
	cases := []struct {
		name       string
		tripID     string
		accountID  string
		tripStatus rentalpb.TripStatus
		wantErr    bool
	}{
		{
			name:       "finished",
			tripID:     "61358659c66699ba60b09754",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		}, {
			name:       "annother_finished",
			tripID:     "61358659c66699ba60b09755",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		}, {
			name:       "in_progress",
			tripID:     "61358659c66699ba60b09756",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		}, {
			name:       "another_in_progress",
			tripID:     "61358659c66699ba60b09757",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    true,
		}, {
			name:       "in_progress_by_another_account",
			tripID:     "61358659c66699ba60b09758",
			accountID:  "account2",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	for _, cc := range cases {
		mgutil.NewObjIDWithValue(objid.MustFromID(id.TripID(cc.tripID)))
		tr, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: cc.accountID,
			Status:    cc.tripStatus,
		})
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: error expected; got none", cc.name)
			}
			continue
		}

		if err != nil {
			t.Errorf("%s: error creating trip: %v", cc.name, err)
			continue
		}

		if tr.ID.Hex() != cc.tripID {
			t.Errorf("%s: incorrect trip id; want: %q; got:%q", cc.name,
				cc.tripID, tr.ID.Hex())
		}

	}

}
func TestCetTrip(t *testing.T) {
	//start container
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongdo db:%v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	acct := id.AccountID("account1")
	mgutil.NewObjID = primitive.NewObjectID
	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: acct.String(),
		CarId:     "car1",
		Start: &rentalpb.LocationStatus{
			PoiName: "startPoint",
			Location: &rentalpb.Location{
				Latitude:  30,
				Longitude: 210,
			},
		},
		End: &rentalpb.LocationStatus{
			PoiName:  "endpoint",
			FeeCent:  10000,
			KmDriven: 35,
			Location: &rentalpb.Location{
				Latitude:  35,
				Longitude: 115,
			},
		}, Status: rentalpb.TripStatus_FINISHED,
	})
	if err != nil {
		t.Fatalf("cannot create trip :%v", err)
	}

	got, err := m.GetTrip(c, objid.ToTripID(tr.ID), acct)
	if err != nil {
		t.Errorf("cannot get trip: %v", err)
	}
	if diff := cmp.Diff(tr, got, protocmp.Transform()); diff != "" {
		t.Errorf("result differs; -want +got: %s", diff)
	}
}

func TestCetTrips(t *testing.T) {
	rows := []struct {
		id        string
		accountID id.AccountID
		status    rentalpb.TripStatus
	}{
		{
			id:        "61358659c66699ba60b09760",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "61358659c66699ba60b09761",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "61358659c66699ba60b09762",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "61358659c66699ba60b09763",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			id:        "61358659c66699ba60b09764",
			accountID: "account_id_for_get_trips1",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongdo db:%v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	for _, r := range rows {
		mgutil.NewObjIDWithValue(id.TripID(r.id))
		_, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: string(r.accountID),
			Status:    r.status,
		})
		if err != nil {
			t.Fatalf("cannot create rows: %v", err)
		}

	}
	cases := []struct {
		name      string
		accountID string
		status    rentalpb.TripStatus
		wantCount int
		wanOnlyID string
	}{
		{
			name:      "get_all",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_TS_NOT_SPECIFIED,
			wantCount: 4,
		},
		{
			name:      "get_in_progress",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_IN_PROGRESS,
			wantCount: 1,
			wanOnlyID: "61358659c66699ba60b09763",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			res, err := m.GetTrips(context.Background(), id.AccountID(cc.accountID), cc.status)
			if err != nil {
				t.Errorf("cannot get trips: %v", err)
			}
			if cc.wantCount != len(res) {
				t.Errorf("incorrect result count; want %d, got %d", cc.wantCount, len(res))
			}
			if cc.wanOnlyID != "" && len(res) > 0 {
				if cc.wanOnlyID != res[0].ID.Hex() {
					t.Errorf("only_id incorrect; want %q got %q", cc.wanOnlyID, res[0].ID.Hex())
				}
			}
		})
	}
}

func TestUpdateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongdo db:%v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	var now int64 = 10000
	tid := id.TripID("61358659c66699ba60b09766")
	aid := id.AccountID("account_for_update")
	mgutil.NewObjIDWithValue(tid)
	mgutil.UpdatedAt = func() int64 {
		return now
	}

	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi",
		},
	})
	if err != nil {
		t.Fatalf("cannot create trip: %v", err)
	}
	if tr.UpdateAt != 10000 {
		t.Fatalf("wrong updateat; want 1000, got: %d", tr.UpdateAt)
	}

	update := &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi_updated",
		},
	}
	cases := []struct {
		name          string
		now           int64
		withUpdatedAt int64
		wantErr       bool
	}{
		{
			name:          "normal_update",
			now:           20000,
			withUpdatedAt: 10000,
			wantErr:       false,
		},
		{
			name:          "update_with_stale_timestamp",
			now:           30000,
			withUpdatedAt: 10000,
			wantErr:       true,
		},
		{
			name:          "update_with_refetch_timestamp",
			now:           40000,
			withUpdatedAt: 20000,
			wantErr:       false,
		},
	}

	for _, cc := range cases {
		now = cc.now
		err := m.UpdateTrip(c, tid, aid, cc.withUpdatedAt, update)
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want error ; got none", cc.name)
			} else {
				continue
			}
		} else {
			if err != nil {
				t.Errorf("%s: cannot update: %v", cc.name, err)
			}
		}
		updatedTrip, err := m.GetTrip(c, tid, aid)
		if err != nil {
			t.Errorf("%s: cannot get trip after update: %v", cc.name, err)
		}
		if cc.now != updatedTrip.UpdateAt {
			t.Errorf("%s: incorrect updateat: want %d, got %d", cc.name, cc.now, updatedTrip.UpdateAt)
		}

	}

}
func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
