package main

import (
	"github.com/konstfish/aquarium/common/db"
	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/konstfish/aquarium/sprite/mappings"
)

func main() {
	monitoring.InitTracer("sprite")

	db.InitRedis()

	mappings.CreateUrlMappings("./sprites")
	mappings.Router.Run(":4001")
}
