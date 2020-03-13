package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"log"
	"net/http"

	"github.com/iissy/go-micro/config"
	"github.com/iissy/go-micro/helloworld"
	"github.com/iissy/go-micro/messages"

	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-plugins/registry/consul/v2"

	wrapperPrometheus "github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Greeter struct{}

func (s *Greeter) SayHello(ctx context.Context, req *messages.HelloRequest, rsp *messages.HelloReply) error {
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
		micro.WrapHandler(wrapperPrometheus.NewHandlerWrapper()),
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
	http.Handle("/metrics", promhttp.Handler())
	// 启动web服务，监听8085端口
	go func() {
		err := http.ListenAndServe("192.168.111.1:8085", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}