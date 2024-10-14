package main

import (
	"github.com/nrgx/wow-pow-tcp/internal/client"
	"github.com/nrgx/wow-pow-tcp/internal/config"
)

func main() {
	conf := config.New()

	client.Run(conf.Addr())
}
