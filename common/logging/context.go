package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// funny mutex

type ContextLogger struct {
	*log.Logger
	mu sync.Mutex
}

var (
	globalLogger *ContextLogger
	once         sync.Once
)

func GetLogger() *ContextLogger {
	once.Do(func() {
		globalLogger = &ContextLogger{
			Logger: log.New(os.Stdout, "", 0),
		}
	})
	return globalLogger
}

func (cl *ContextLogger) log(ctx context.Context, level string, messages ...string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	traceID := extractTrace(ctx)
	combinedMessage := strings.Join(messages, " ")

	fmt.Printf("[COM] %s | %s | TraceID: %v | %s\n",
		time.Now().Format("2006/01/02 - 15:04:05"),
		level,
		traceID,
		combinedMessage,
	)
}

func (cl *ContextLogger) Info(ctx context.Context, messages ...string) {
	cl.log(ctx, "INF", messages...)
}

func (cl *ContextLogger) Error(ctx context.Context, messages ...string) {
	cl.log(ctx, "ERR", messages...)
}

func Info(ctx context.Context, messages ...string) {
	GetLogger().Info(ctx, messages...)
}

func Error(ctx context.Context, messages ...string) {
	GetLogger().Error(ctx, messages...)
}
