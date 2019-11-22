package main

import (
	"context"
	"log"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-plugins/registry/consul"
	hello "greeter/src/helloworld"
)

// Greeter is
type Greeter struct{}

// SayHello is
func (s *Greeter) SayHello(ctx context.Context, req *hello.HelloRequest, rsp *hello.HelloReply) error {
	rsp.Message = "Hello " + req.Name
	return nil
}

func main() {
	// 修改consul地址，如果是本机，这段代码和后面的那行使用代码都是可以不用的
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"192.168.111.148:8500",
		}
	})

	service := grpc.NewService(
		micro.Registry(reg),
		micro.Name("greeter"),
	)

	service.Init()

	hello.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
