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
		mgutil.NewObjID = func() primitive.ObjectID {
			return objid.MustFromID(id.TripID(cc.tripID))
		}
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

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
