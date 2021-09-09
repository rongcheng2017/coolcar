package mongotesting

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	image         = "mongo:4.4"
	containerPort = "27017/tcp"
)

var mongoURI string

const defaultMongoURI = "mongdb://localhost:27017"

// RunWithMongoInDocker runs the tests with
// a mongodb instance in a docker container.
func RunWithMongoInDocker(m *testing.M) int {
	//本机的docker环境
	c, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	//创建容器
	resp, err := c.ContainerCreate(ctx,
		&container.Config{
			Image: image,
			ExposedPorts: nat.PortSet{
				containerPort: {},
			},
		}, &container.HostConfig{
			PortBindings: nat.PortMap{
				containerPort: []nat.PortBinding{
					{
						HostIP:   "127.0.0.1",
						HostPort: "0", //0自动挑选端口
					},
				},
			},
		}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	containerID := resp.ID
	//移除
	defer func() {
		err = c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err)
		}
	}()
	//启动
	err = c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	//获取端口
	inspRes, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}
	hostPort := inspRes.NetworkSettings.Ports[containerPort][0]
	mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)

	return m.Run()
}

func NewClient(c context.Context) (*mongo.Client, error) {
	if mongoURI == "" {
		return nil, fmt.Errorf("mong uri not set. Please run RunWithMongoInDocker in TestMain ")
	}
	return mongo.Connect(c, options.Client().ApplyURI(mongoURI))
}

func NewDefaultCleint(c context.Context) (*mongo.Client, error) {
	return mongo.Connect(c, options.Client().ApplyURI(defaultMongoURI))
}

//SetupIndexes sets up indexes for the given database.
func SetupIndexes(c context.Context, d *mongo.Database) error {
	_, err := d.Collection("account").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{
				Key: "open_id", Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	_, err = d.Collection("trip").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{Key: "trip.accountid", Value: 1},
			{Key: "trip.status", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
			"trip.status": 1,
		}),
	})
	if err != nil {
		return err
	}

	_, err = d.Collection("profile").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{Key: "accountid", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	return err

}
