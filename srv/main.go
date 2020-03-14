package main

import (
	"context"
	"github.com/iissy/go-micro/config"
	"github.com/iissy/go-micro/helloworld"
	"github.com/iissy/go-micro/messages"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"

	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-plugins/registry/consul/v2"

	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
)

type Greeter struct{}

func (s *Greeter) SayHello(ctx context.Context, req *messages.HelloRequest, rsp *messages.HelloReply) error {
	log.Printf("Received: %s", req.Name)
	rsp.Message = "Hello, " + req.Name
	return nil
}

func main() {
	urls := config.GetConsulUrls()

	// 修改consul地址，如果是本机，这段代码和后面的那行使用代码都是可以不用的
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = urls
	})

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.srv.greeter"),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// 优雅关闭服务
	service.Server().Init(
		server.Wait(nil),
	)

	// 如果单机启动多个实例要启动不同的端口8085,8086...
	prometheusBoot()
	helloworld.RegisterGreeterHandler(service.Server(), new(Greeter))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func prometheusBoot() {
	service := web.NewService(
		web.Name("go.micro.srv.metrics"),
	)
	service.Handle("/metrics", promhttp.Handler())

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := service.Run(); err != nil {
			log.Fatal(err)
		}
	}()
}
