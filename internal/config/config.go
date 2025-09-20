package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env" env-default:"development"`
	StoragePath    string `yaml:"storage_path" env-required:"true"`
	MigrationPath  string `yaml:"migration_path"`
	MigrationTable string `yaml:"migration_table"`
	HTTPServer     `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

// Функция чтения config-файла, заполнения ее значениями структуры Config
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %v", err)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}
	return &cfg
}
