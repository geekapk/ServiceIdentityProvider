package SIPConfig

import (
	"encoding/json"
)

type Config struct {
	ListenAddr string
	Services []*Service
	Groups []*Group
}

type Service struct {
	Name string
	Key string
}

type Group struct {
	Name string
	Key string
	Members []string
}

func LoadConfig(data []byte) (*Config, error) {
	config := &Config {}
	err := json.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
