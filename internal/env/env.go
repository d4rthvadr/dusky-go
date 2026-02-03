package env

import (
	"os"
	"strconv"
)

func GetEnv(key, defaultValue string) string {

	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func GetEnvAsInt(key string, defaultValue int) int {

	valStr, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	valAsInt, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return valAsInt
}
