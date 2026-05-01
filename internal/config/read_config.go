package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return homeDir, err
	}

	filePath := homeDir + "/" + configFileName

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fileInfo.Name(), fmt.Errorf("~/%v not found!", configFileName)
	}

	return filePath, nil
}

func ReadGatorconfigContent() {
	filePath, err := getConfigFilePath()
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}

func Read() Config {

	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}
	}

	gatorconfig, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return Config{}
	}

	cfg := Config{}
	err = json.Unmarshal(gatorconfig, &cfg)
	if err != nil {
		return cfg
	}

	return cfg
}
