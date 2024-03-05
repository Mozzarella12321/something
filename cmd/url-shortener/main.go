package main

import (
	"fmt"

	"github.com/mozzarella12321/orders-api/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg.StoragePath) //delete

	//TODO: init logger: slog

	//TODO: init storage: sqlite

	//TODO: init router: chi, "chi render"

	//TODO: run server:
}
