package main

import (
	"fmt"
	"github.com/BioMihanoid/URLShortener/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
