package mq

import (
	"context"
	"coolcar/car/api/gen/v1"

)

type Subscriber interface {
	Subscribe(context.Context) (chan *carpb.CarEntity, func(), error)
}
type Publisher interface {
	Publish(context.Context, *carpb.CarEntity) error
}

