package config

type Config struct {
	DataBasePath string
	Port         string
}

func New() Config {
	return Config{
		DataBasePath: "./internal/database/social.db",
		Port:         "8000",
	}
}
