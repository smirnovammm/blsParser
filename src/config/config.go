package config

import (
	"encoding/json"
	"io/ioutil"
)

const configFile = "./config.json"

type Email struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type Config struct {
	Email   Email  `json:"email"`
	TgToken string `json:"tg_token"`
}

func Load() (*Config, error) {
	var config Config
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
