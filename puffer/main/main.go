package main

import (
	"github.com/konstfish/aquarium/common/db"
	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/konstfish/aquarium/puffer/mappings"
)

func main() {
	monitoring.InitTracer("puffer")

	db.InitRedis()

	mappings.CreateUrlMappings()
	mappings.Router.Run(":4003")
}
