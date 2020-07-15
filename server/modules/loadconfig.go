package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
)

var (
	Config Configuration
)

type Configuration struct {
	Logging struct {
		LogFile   string
		UploadDir string
	}
	TCP struct {
		Enabled bool
		Port    int
	}
	UDP struct {
		Enabled bool
		Port    int
	}
	HTTP struct {
		Enabled    bool
		Port       int
		Get        string
		Post       string
		Uploadsize int64
	}
	HTTPS struct {
		Enabled     bool
		Port        int
		Get         string
		Post        string
		Uploadsize  int64
		Certificate string
		Key         string
	}
	ICMP struct {
		Enabled bool
	}
	DNS struct {
		Enabled bool
	}
}

func LoadConf(config *string) (Configuration, error) {
	// Make sure the configuration file is a valid file
	if !IsValidFile(*config) {
		return Config, errors.New("Configuration file does not exist: " + *config)
	}

	file, err := os.Open(*config)
	if err != nil {
		fmt.Println("Could not open configuration file: " + *config)
		return Config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	//Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Could not decode configuration file: " + *config)
		return Config, err
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	Config.Logging.UploadDir = path.Join(pwd, Config.Logging.UploadDir)
	Config.Logging.LogFile = path.Join(pwd, Config.Logging.LogFile)

	fmt.Println("Configuration file loaded successfully")
	return Config, nil
}
