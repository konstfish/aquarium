package main

import (
	"github.com/konstfish/aquarium/common/db"
	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/konstfish/aquarium/tetra/mappings"
)

func main() {
	monitoring.InitTracer("tetra")

	db.InitRedis()

	mappings.CreateUrlMappings()
	mappings.Router.Run(":4002")
}
