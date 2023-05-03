package config

import (
	"encoding/json"
	"log"
	"os"
)

type server struct {
	Addr string `json:"addr"`
}

type database struct {
	Driver string `json:"driver"`
	Dsn    string `json:"dsn"`
}

type Config struct {
	Server   server   `json:"server"`
	Database database `json:"database"`
}

func NewConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("read file in config %v", err)
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("unmarshal in config %v", err)
		return nil, err
	}
	return &config, nil
}
