package dao

import (
	"context"
	"coolcar/shared/id"
	"coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func TestResolveAccountID(t *testing.T) {
	//start container
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongdo db:%v", err)
	}
	m := NewMongo(mc.Database("coolcar"))
	//init db
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("612cb3cedd1930deb67c9a8e")),
			openIDField:     "openid_1",
		},
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("612cb3cedd1930deb67c9a70")),
			openIDField:     "openid_2",
		},
	})
	if err != nil {
		t.Fatalf("cannot insert initial values: %v", err)
	}
	mgutil.NewObjID = func() primitive.ObjectID {
		return objid.MustFromID(id.AccountID("612cb3cedd1930deb67c9a71"))
	}

	//test case
	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "openid_1",
			want:   "612cb3cedd1930deb67c9a8e",
		}, {
			name:   "another_existing_user",
			openID: "openid_2",
			want:   "612cb3cedd1930deb67c9a70",
		},
		{
			name:   "new_user",
			openID: "openid_3",
			want:   "612cb3cedd1930deb67c9a71",
		},
	}

	//test
	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), cc.openID)
			if err != nil {
				t.Errorf("faild resolve account id for %q:%v", cc.openID, err)
			} else {
				want := cc.want
				if id.String() != want {
					t.Errorf("resulve account id : want:%q,got:%q", want, id)
				}
			}
		})
	}

}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
