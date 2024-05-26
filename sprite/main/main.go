package main

import (
	"github.com/konstfish/aquarium/sprite/mappings"
)

func main() {
	mappings.CreateUrlMappings("./sprites")
	mappings.Router.Run(":4001")
}
