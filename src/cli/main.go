package main

import (
	"fmt"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-plugins/registry/consul"
	hello "greeter/src/helloworld"

	"context"
)

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

	// use the generated client stub
	cl := hello.NewGreeterService("greeter", service.Client())

	// Set arbitrary headers in context
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"X-User-Id": "john",
		"X-From-Id": "script",
	})

	rsp, err := cl.SayHello(ctx, &hello.HelloRequest{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Message)
}
