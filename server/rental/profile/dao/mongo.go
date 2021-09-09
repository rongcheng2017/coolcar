package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	accountIDField      = "accountid"
	profileField        = "profile"
	identityStatusField = profileField + ".identitystatus"
)

type Mongo struct {
	col *mongo.Collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("profile"),
	}
}

type ProfileRecord struct {
	AccountID string            `bson:"accountid"`
	Profile   *rentalpb.Profile `bson:"profile"`
}

func (m *Mongo) GetProfile(c context.Context, aid id.AccountID) (*rentalpb.Profile, error) {
	res := m.col.FindOne(c, byAccountID(aid))
	if err := res.Err(); err != nil {
		return nil, err
	}
	var pr ProfileRecord
	err := res.Decode(&pr)
	if err != nil {
		return nil, fmt.Errorf("cannot decode profile record:%v", err)
	}
	return pr.Profile, nil
}

func (m *Mongo) UpdateProfile(c context.Context, aid id.AccountID, preStatus rentalpb.IdentityStatus, profile *rentalpb.Profile) error {
	_, err := m.col.UpdateOne(c, bson.M{
		accountIDField:      aid.String(),
		identityStatusField: preStatus,
	}, mgutil.Set(bson.M{
		profileField:   profile,
		accountIDField: aid.String(),
	}), options.Update().SetUpsert(true))
	return err

}

func byAccountID(aid id.AccountID) bson.M {
	return bson.M{
		accountIDField: aid.String()}
}
