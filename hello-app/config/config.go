package config

import (
	"github.com/junydania/pppro/hello-app/utils/env"
	"strconv"
)

type Config struct {
	Port        int
	Timeout     int
	Region string
	Environment string
	DB_HOST string
	DB_PORT string
}


func GetConfig() Config {
	return Config{
		Port:        parseEnvToInt("PORT", "8082"),
		Timeout:     parseEnvToInt("TIMEOUT", "30"),
		Region: env.GetEnv("REGION", "ap-southeast-1"),
		Environment: env.GetEnv("ENV", "development"),
		DB_HOST: env.GetEnv("DB_HOST", "localhost"),
		DB_PORT: env.GetEnv("DB_PORT", "8000"),
	}
}

func parseEnvToInt(envName, defaultValue string) int {
	num, err := strconv.Atoi(env.GetEnv(envName, defaultValue))

	if err != nil {
		return 0
	}
	return num
}
