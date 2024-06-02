package main

import (
	"github.com/konstfish/aquarium/butterfly/mappings"
	"github.com/konstfish/aquarium/common/db"
	"github.com/konstfish/aquarium/common/monitoring"
)

func main() {
	monitoring.InitTracer("butterfly")

	db.InitRedis()

	mappings.CreateUrlMappings()
	mappings.Router.Run(":4004")
}
