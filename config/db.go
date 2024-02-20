package config

import "os"

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func GetDBConfig() DBConfig {
	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		user = "brick"
	}

	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		password = "brick"
	}

	database, ok := os.LookupEnv("DB_DATABASE")
	if !ok {
		database = "brick"
	}

	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		host = "localhost"
	}

	port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		port = "5432"
	}

	return DBConfig{
		User:     user,
		Password: password,
		Database: database,
		Host:     host,
		Port:     port,
	}

}
