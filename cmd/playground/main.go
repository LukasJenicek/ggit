package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/LukasJenicek/ggit/internal/config"
)

func main() {
	cfg, err := config.LoadGitConfig()
	if err != nil {
		panic(err)
	}

	spew.Dump(cfg)
}
