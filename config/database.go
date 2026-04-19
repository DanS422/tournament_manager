package config

import "os"

type Config struct {
	Env            string
	DatabaseURL    string
	MigrationsPath string
}

func Load() Config {
	cfg := Config{
		Env:         getEnv("APP_ENV", "development"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = defaultDB(cfg.Env)
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func defaultDB(env string) string {
	switch env {
	case "test":
		return "sqlite://file:test.db?mode=memory&cache=shared"
	case "production":
		return "sqlite://tournaments.db"
	default:
		return "sqlite://tournaments.db"
	}
}
