package main

import (
	"fmt"
	"github.com/iissy/go-micro/helloworld"
	"github.com/iissy/go-micro/messages"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"

	"context"
)

func main() {
	// 修改consul地址，如果是本机，这段代码和后面的那行使用代码都是可以不用的
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"47.244.143.251:8500",
		}
	})

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("greeter"),
	)

	service.Init()

	// use the generated client stub
	cl := helloworld.NewGreeterService("greeter", service.Client())

	// Set arbitrary headers in context
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"X-User-Id": "john",
		"X-From-Id": "script",
	})

	rsp, err := cl.SayHello(ctx, &messages.HelloRequest{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.Message)
}
