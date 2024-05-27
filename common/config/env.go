package config

import (
	"os"
)

func GetConfigVar(variable string) string {
	return os.Getenv(variable)
}
