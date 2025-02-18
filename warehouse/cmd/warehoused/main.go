package main

import (
	"github.com/Jerry19900615/go-shopping/shipping/proto"
	"github.com/Jerry19900615/go-shopping/warehouse/internal/platform/broker"
	"github.com/Jerry19900615/go-shopping/warehouse/internal/platform/config"
	"github.com/Jerry19900615/go-shopping/warehouse/internal/platform/redis"
	"github.com/Jerry19900615/go-shopping/warehouse/internal/service"
	"github.com/Jerry19900615/go-shopping/warehouse/proto"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	gmbroker "github.com/micro/go-micro/broker"
	"time"
)

func main() {

	if err := gmbroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := gmbroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}
	itemShippedChannel := make(chan *shipping.ItemShippedEvent)
	broker.CreateEventConsumer(itemShippedChannel)

	repo := redis.NewWarehouseRepository(":6379")

	svc := grpc.NewService(
		micro.Name(config.ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(config.Version),
	)
	svc.Init()

	warehouse.RegisterWarehouseHandler(svc.Server(), service.NewWarehouseService(repo, itemShippedChannel))

	if err := svc.Run(); err != nil {
		panic(err)
	}
}
