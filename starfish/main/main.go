package main

import (
	"github.com/konstfish/aquarium/common/db"
	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/konstfish/aquarium/starfish/controllers"
	"github.com/konstfish/aquarium/starfish/mappings"
)

func main() {
	monitoring.InitTracer("starfish")

	db.InitRedis()

	redisClient := db.ConnectRedis()
	defer redisClient.Client.Close()
	go redisClient.ListenForNewItems("starfish", controllers.ProcessEventAmount)

	mappings.CreateUrlMappings()
	mappings.Router.Run(":4005")
}
