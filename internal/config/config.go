package config

import (
    "os"
    "gopkg.in/yaml.v3"
)

type BackendConfig struct {
    Address string `yaml:"address"`
}

type Config struct {
    Port                int             `yaml:"port"`
    HealthCheckInterval int             `yaml:"health_check_interval"`
    Backends            []BackendConfig `yaml:"backends"`
}

func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    err = yaml.Unmarshal(data, &cfg)
    if err != nil {
        return nil, err
    }

    return &cfg, nil
}