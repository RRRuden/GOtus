package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type HTTPServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type BookingServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	Env           string              `yaml:"env"`
	StoragePath   string              `yaml:"storage_path"`
	HTTPServer    HTTPServerConfig    `yaml:"http_server"`
	BookingServer BookingServerConfig `yaml:"booking_server"`
	MongoDB       MongoDBConfig       `yaml:"mongodb"`
	Redis         RedisConfig         `yaml:"redis"`
	PostgreSQL    PostgresConfig      `yaml:"postgresql"`
}

func LoadConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Ошибка открытия конфигурационного файла: %v", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("Ошибка декодирования конфигурации: %v", err)
	}

	return &cfg
}
