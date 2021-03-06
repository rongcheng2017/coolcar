package poi

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"hash/fnv"

	"github.com/gogo/protobuf/proto"
)

var poi = []string{
	"中关村",
	"天安门",
	"陆家嘴",
	"迪士尼",
	"天河体育中心",
	"广州塔",
}

type Manager struct {
}

func (pm *Manager) Resolve(c context.Context, ls *rentalpb.Location) (string, error) {
	b, err := proto.Marshal(ls)
	if err != nil {
		return "", err
	}
	h := fnv.New32()
	h.Write(b)
	return poi[int(h.Sum32())%len(poi)], nil
}
