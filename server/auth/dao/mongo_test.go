package dao

import (
	"context"
	mgo "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI string

func TestResolveAccountID(t *testing.T) {
	//start container
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("cannot connect mongdo db:%v", err)
	}
	m := NewMongo(mc.Database("coolcar"))
	//init db 
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgo.IDField: mustObjecId("612cb3cedd1930deb67c9a8e"),
			openIDField: "openid_1",
		},
		bson.M{
			mgo.IDField: mustObjecId("612cb3cedd1930deb67c9a70"),
			openIDField: "openid_2",
		},
	})
	if err != nil {
		t.Fatalf("cannot insert initial values: %v", err)
	}
	m.newObjID = func() primitive.ObjectID {
		return mustObjecId("612cb3cedd1930deb67c9a71")
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
				t.Errorf("faild resolve account id for %q:%v",cc.openID, err)
			} else {
				want := cc.want
				if id != want {
					t.Errorf("resulve account id : want:%q,got:%q", want, id)
				}
			}
		})
	}

}
func mustObjecId(hex string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return objID
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
}
