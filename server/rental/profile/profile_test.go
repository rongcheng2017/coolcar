package profile

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/profile/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	mongotesting "coolcar/shared/mongo/testing"
	"coolcar/shared/server"
	"os"
	"testing"

)

func TestProfileLifecycler(t *testing.T) {

	c := context.Background()
	s := newService(c,t)

	aid := id.AccountID("account1")
	c = auth.ContextWithAccountID(c, aid)

	cases := []struct {
		name       string
		op         func() (*rentalpb.Profile, error)
		wantName   string
		wantStatus rentalpb.IdentityStatus
		wantErr    bool
	}{
		{
			name: "get_empty",
			op: func() (*rentalpb.Profile, error) {
				return s.GetProfile(c, &rentalpb.GetProfileRequest{})
			},
			wantStatus: rentalpb.IdentityStatus_UNSUBMITTED,
		},
		{
			name: "submit",
			op: func() (*rentalpb.Profile, error) {
				return s.SubmitProfile(c, &rentalpb.Identity{
					Name: "abc",
				})
			},
			wantName:   "abc",
			wantStatus: rentalpb.IdentityStatus_PENDING,
		}, {
			name: "submit_again",
			op: func() (*rentalpb.Profile, error) {
				return s.SubmitProfile(c, &rentalpb.Identity{
					Name: "abc",
				})
			},
			wantErr: true,
		},
		{
			name: "todo_force_verify",
			op: func() (*rentalpb.Profile, error) {
				p := &rentalpb.Profile{
					Identity: &rentalpb.Identity{
						Name: "abc",
					},
					IdentityStatus: rentalpb.IdentityStatus_VERIFIED,
				}
				err := s.Mongo.UpdateProfile(c, aid, rentalpb.IdentityStatus_PENDING, p)
				if err != nil {
					return nil, err
				}
				return p, nil
			},
			wantName:   "abc",
			wantStatus: rentalpb.IdentityStatus_VERIFIED,
		}, {
			name: "clear",
			op: func() (*rentalpb.Profile, error) {
				return s.ClearProfile(c, &rentalpb.ClearProfileRequest{})
			},
			wantStatus: rentalpb.IdentityStatus_UNSUBMITTED,
		},
	}
	for _, cc := range cases {
		p, err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want error: got none", cc.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s: opration failure :%v",cc.name,err)
		}
		gotName := ""
		if p.Identity!=nil {
			gotName=p.Identity.Name
		}
		if gotName!=cc.wantName {
			t.Errorf("%s: want name is %s, got is %s", cc.name, cc.wantName, p.Identity.Name)
		}
		if p.IdentityStatus != cc.wantStatus {
			t.Errorf("%s: want status is %v, got is %v", cc.name, cc.wantStatus, p.IdentityStatus)
		}
	}
}
func newService(c context.Context, t *testing.T) *Service {
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot create new mongo client: %v", err)
	}
	db := mc.Database("profile")
	mongotesting.SetupIndexes(c, db)
	logger, err := server.NewZapLogger()
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}
	return &Service{
		Mongo:  dao.NewMongo(db),
		Logger: logger,
	}
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}