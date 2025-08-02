package main

import (
	"log"

	"github.com/sirawatc/simple-gin-crud/internal/shared/config"
)

func main() {
	cfg := config.NewConfig()
	log.Println("Initializing completed with config: ", cfg)
}
