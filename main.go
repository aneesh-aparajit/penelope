package main

import (
	"fmt"
	"log/slog"

	"github.com/aneesh-aparajit/penelope/pkg/config"
)

func main() {
	fmt.Println("Hello World")
	cfg, err := config.LoadConfig(".penelope.yml")
	if err != nil {
		slog.Error("failed to read config", slog.Any("err", err))
		panic(err)
	}
	fmt.Println(cfg)
}
