package main

import (
	"fmt"
	"loadbalancer/server"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Backends []string `yaml:"backends"`
	Port     int      `yaml:"port"`
}

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		fmt.Println("Error parsing YAML:", err)
		return
	}
	tcpServer := server.NewServer(c.Backends, c.Port)
	tcpServer.Listen()
}
