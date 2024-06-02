package db

import (
	"context"
	"encoding/json"
	"log"

	"github.com/konstfish/aquarium/common/config"
	"github.com/konstfish/aquarium/common/monitoring"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var Redis *RedisClient

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

type RedisQueueItem struct {
	Data        string `json:"data"`
	TraceParent string `json:"traceparent"`
}

func (item *RedisQueueItem) Serialize() (error, string) {
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err, ""
	}

	return nil, string(itemJSON)
}

func (item *RedisQueueItem) Deserialize(itemJSON string) error {
	err := json.Unmarshal([]byte(itemJSON), &item)
	if err != nil {
		return err
	}

	return nil
}

func InitRedis() {
	Redis = ConnectRedis()
}

func ConnectRedis() *RedisClient {
	opt, err := redis.ParseURL(config.GetConfigVar("REDIS_URI"))
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Println(err)
	}

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func (r *RedisClient) ListenForNewItems(queueName string, handler func(ctx context.Context, msg string)) {
	log.Println("Listening for new items in queue", queueName)

	for {
		var ctx context.Context
		var span trace.Span

		ctx = context.Background()

		// pop item from queue
		result, err := r.Client.BLPop(ctx, 0, queueName).Result()
		if err != nil {
			log.Println(err)
		}

		// deserialize queue item
		var queueItem RedisQueueItem
		err = queueItem.Deserialize(result[1])
		if err != nil {
			log.Println(err)
		}

		// create span
		sc, err := monitoring.ParseTraceparentHeader(queueItem.TraceParent)
		if err == nil {
			ctx, span = monitoring.Tracer.Start(
				trace.ContextWithRemoteSpanContext(ctx, sc),
				(queueName + " receive"),
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(
					attribute.String("messaging.system", "redis"),
					attribute.String("messaging.operation", "receive"),
					attribute.String("messaging.destination.name", queueName),
				),
			)
		}

		handler(ctx, queueItem.Data)

		if span != nil {
			span.End()
		}
	}
}

func (r *RedisClient) PushToQueue(ctx context.Context, queueName string, value string) {
	log.Printf("Pushing %s to queue %s", value, queueName)

	var traceparent = monitoring.EmptyTraceparentHeader()

	if monitoring.Tracer != nil {
		var span trace.Span

		ctx, span = monitoring.Tracer.Start(
			ctx,
			(queueName + " publish"),
			trace.WithSpanKind(trace.SpanKindProducer),
			trace.WithAttributes(
				attribute.String("messaging.system", "redis"),
				attribute.String("messaging.operation", "publish"),
				attribute.String("messaging.destination.name", queueName),
			),
		)
		defer span.End()

		traceparent = monitoring.ExtractTraceparentHeader(ctx)
	}

	// create queue item
	queueItem := RedisQueueItem{
		Data:        value,
		TraceParent: traceparent,
	}

	// serialize queue item
	err, item := queueItem.Serialize()
	if err != nil {
		log.Println(err)
	}

	r.Client.RPush(ctx, queueName, item)
}

func (r *RedisClient) PushToQueueWithDefaultContext(queueName string, value string) {
	r.PushToQueue(r.Ctx, queueName, value)
}
