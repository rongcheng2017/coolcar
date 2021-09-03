package mgutil

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	IDFieldName        = "_id"
	UpdatedAtFieldName = "updatedat"
)

func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}

var NewObjID = primitive.NewObjectID

var UpdatedAt = func() int64 {
	return time.Now().UnixNano()
}

func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}

type IDField struct {
	ID primitive.ObjectID `bson:"_id"`
}
type UpdatedAtField struct {
	UpdateAt int64 `bson:"updatedat"`
}
