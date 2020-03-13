package main

import (
	"context"
	"github.com/iissy/go-micro/config"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iissy/go-micro/helloworld"
	"github.com/iissy/go-micro/messages"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/web"
)

type Say struct{}

var (
	cl helloworld.GreeterService
)

func (s *Say) Anything(c *gin.Context) {
	log.Print("Received Say.Anything API request")
	c.JSON(200, map[string]string{
		"message": "Hi, this is the Greeter API",
	})
}

func (s *Say) Hello(c *gin.Context) {
	log.Print("Received Say.Hello API request")

	name := c.Param("name")

	response, err := cl.SayHello(context.TODO(), &messages.HelloRequest{
		Name: name,
	})

	if err != nil {
		c.JSON(500, err)
	}

	c.JSON(200, response)
}

func main() {
	urls := config.GetConsulUrls()

	// 修改consul地址，如果是本机，这段代码和后面的那行使用代码都是可以不用的
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = urls
	})

	// Create service
	service := web.NewService(
		web.Registry(reg),
		web.Name("go.micro.api.greeter"),
	)

	service.Init()

	// setup Greeter Server Client
	cl = helloworld.NewGreeterService("go.micro.srv.greeter", client.DefaultClient)

	// Create RESTful handler (using Gin)
	say := new(Say)
	router := gin.Default()
	router.GET("/greeter", say.Anything)
	router.GET("/greeter/:name", say.Hello)

	// Register Handler
	service.Handle("/", router)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
