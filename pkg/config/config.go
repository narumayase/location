package config

type Config struct {
	DBPath string
}

func New(env string) *Config {

	config := &Config{DBPath: "location.db"}

	if env != "local" {
		config.DBPath = "/usr/bin/location.db"
	}
	return config
}
