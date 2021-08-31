package dao

import (
	"context"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

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
	m.newObjID=func() primitive.ObjectID {
		objID,_:=primitive.ObjectIDFromHex("612cb3cedd1930deb67c9a8f")
		return objID
	}
	id, err := m.ResolveAccountID(c, "123")
	if err != nil {
		t.Errorf("faild resolve account id for 123:%v", err)
	} else {
		want := "612cb3cedd1930deb67c9a8f"
		if id != want {
			t.Errorf("resulve account id : want:%q,got:%q", want, id)
		}
	}

}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
}
