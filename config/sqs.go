package config

import "os"

func GetSQSUrl() string {
	env, ok := os.LookupEnv("QUEUE_URL")
	if !ok {
		return "http://localhost:4566/000000000000/pending-transaction"
	}
	return env
}
