package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACKey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
}

func DefaultConfig() Config {
	return Config{
		Port:     3000,
		Env:      "dev",
		Pepper:   "secret-random-string",
		HMACKey:  "secret-hmac-key",
		Database: DefaultPostgresConfig(),
	}
}

func LoadConfig(configReq bool) Config {

	f, err := os.Open(".config")
	if err != nil {

		if configReq {
			panic(err)
		}

		fmt.Println("Using the default connfig...")
		return DefaultConfig()
	}
	var c Config

	dec := json.NewDecoder(f)

	err = dec.Decode(&c)

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully Loaded .config")

	return c
}
