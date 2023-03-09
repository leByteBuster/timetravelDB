package databaseapi

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Neo4jConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type TimescaleConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func LoadConfig() (Neo4jConfig, TimescaleConfig, error) {

	var config struct {
		Neo4j     Neo4jConfig     `yaml:"neo4j"`
		Timescale TimescaleConfig `yaml:"timescale"`
	}

	b, err := os.ReadFile("config.yaml")
	if err != nil {
		return Neo4jConfig{}, TimescaleConfig{}, err
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	} // err = yaml.Unmarshal(b, &config)
	if err != nil {
		return Neo4jConfig{}, TimescaleConfig{}, err
	}

	return config.Neo4j, config.Timescale, nil
}
