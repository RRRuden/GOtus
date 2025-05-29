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

type Config struct {
	Env         string           `yaml:"env"`
	StoragePath string           `yaml:"storage_path"`
	HTTPServer  HTTPServerConfig `yaml:"http_server"`
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
