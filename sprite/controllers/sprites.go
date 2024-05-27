package controllers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/konstfish/aquarium/common/db"
	"github.com/konstfish/aquarium/common/monitoring"
	"go.opentelemetry.io/otel/trace"
)

func GetRandomSprite(ctx context.Context, rootDir string, spriteDir []string) string {
	return GetSprite(ctx, rootDir+"/"+spriteDir[rand.Intn(len(spriteDir))])
}

func GetSprite(ctx context.Context, spritePath string) string {
	var span trace.Span
	ctx, span = monitoring.Tracer.Start(ctx, "GetSprite")
	defer span.End()

	file, err := db.Redis.Client.Get(ctx, fmt.Sprintf("sprite-%s", spritePath)).Result()
	if err != nil {
		log.Println("cache miss")
		file, err := getSpriteFile(ctx, spritePath)
		if err != nil {
			return ":("
		}

		return file
	}

	return file
}

func getSpriteFile(ctx context.Context, spritePath string) (string, error) {
	var span trace.Span
	ctx, span = monitoring.Tracer.Start(ctx, "getSpriteFile")
	defer span.End()

	// add artificial delay on text load
	time.Sleep(time.Millisecond * 250)

	file, err := os.ReadFile(spritePath)
	if err != nil {
		return "", err
	}

	writeSpriteCache(ctx, spritePath, string(file))

	return string(file), nil
}

func writeSpriteCache(ctx context.Context, spritePath string, file string) {
	db.Redis.Client.Set(ctx, fmt.Sprintf("sprite-%s", spritePath), file, time.Minute*1)
}
