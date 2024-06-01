package main

import (
	"github.com/konstfish/aquarium/butterfly/mappings"
	"github.com/konstfish/aquarium/common/monitoring"
)

func main() {
	monitoring.InitTracer("butterfly")

	mappings.CreateUrlMappings()
	mappings.Router.Run(":4004")
}
