package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Routers []RouterConfig `yaml:"routers"`
	Listen  string         `yaml:"listen"`
}

type RouterConfig struct {
	Route string       `yaml:"route"`
	Hooks []HookConfig `yaml:"hooks"`
}

type HookConfig struct {
	URL     string `yaml:"url"`
	Command string `yaml:"command"`
	Freq    string `yaml:"freq"`
}

func GetConfig() *Config {
	config := &Config{}
	if err := yaml.Unmarshal(getConfigBytes(), config); err != nil {
		log.Fatalln(err.Error())
	}
	return config
}

func getConfigBytes() []byte {
	home, _ := os.UserHomeDir()
	pathList := []string{
		"config.yml",
		home + "/.config/notification/config.yml",
		"/etc/notification/config.yml",
	}
	for _, path := range pathList {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			log.Printf("config file found in %s\n", path)
			return b
		}
	}
	log.Fatalf("error: no config file found in any of %v\n", pathList)
	return nil
}
