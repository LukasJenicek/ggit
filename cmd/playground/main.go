package main

import (
	"github.com/LukasJenicek/ggit/internal/config"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	spew.Dump(cfg)
}
