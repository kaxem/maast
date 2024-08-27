package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Service struct {
		Framework string `yaml:"framework"`
		Port      int    `yaml:"port"`
	} `yaml:"service"`
}

func dockerFileCreate(config Config) error {
	port := config.Service.Port
	lines := []string{
		"FROM python:3.11",
		fmt.Sprintf("ARG PORT=%v", port),
		"RUN mkdir -p /app/src/",
		"WORKDIR /app/src/",
		"ADD . .",
		"RUN pip install -r requirements.txt",
		"EXPOSE $PORT",
		"CMD [\"python\",\"manage.py\",\"runserver\",\"0.0.0.0:$PORT\"]",
	}
	dockerFile, err := os.Create("Dockerfile.maast")
	if err != nil {
		return fmt.Errorf("can't create file %v", err)
	}
	defer dockerFile.Close()

	writer := bufio.NewWriter(dockerFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("couldn't write in dockerfile %v", err)
		}
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing the writer:", err)
	}
	return nil
}

func main() {

	yamlFile, err := ioutil.ReadFile("maast.yaml")
	if err != nil {
		log.Fatalf("error reading yaml file %v", err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("error Unmarshaling yaml")
	}

	dockerFileCreate(config)
}
