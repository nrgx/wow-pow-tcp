package main

import (
	"github.com/nrgx/wow-pow-tcp/internal/config"
	"github.com/nrgx/wow-pow-tcp/internal/server"
)

func main() {
	conf := config.New()

	server.Start(conf.Addr())
}
