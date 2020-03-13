package config

import (
	"fmt"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
	"log"
)

type Address struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

var urls []string

func init() {
	if err := config.Load(file.NewSource(
		file.WithPath("../conf/config.json"),
	)); err != nil {
		log.Panic(err)
	}

	addrs := make([]Address, 0)
	if err := config.Get("consul").Scan(&addrs); err != nil {
		log.Panic(err)
	}

	for _, addr := range addrs {
		urls = append(urls, fmt.Sprintf("%s:%d", addr.Host, addr.Port))
	}
}

func GetConsulUrls() []string {
	return urls
}
