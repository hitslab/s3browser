package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type S3Settings struct {
	BaseUrl    string `json:"base_url"`
	Region     string `json:"region"`
	Bucket     string `json:"bucket"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	PathStyle  bool   `json:"path_style"`
	DisableSSL bool   `json:"disable_ssl"`
}

type Config struct {
	S3 S3Settings `json:"s3"`
}

func getConfigFile() (string, error) {
	filePath, err := xdg.ConfigFile("hitslab/s3browser/config.json")
	if err != nil {
		return "", err
	}

	dir, _ := filepath.Split(filePath)
	if len(dir) == 0 {
		dir = "."
	}

	_, err = os.Stat(dir)

	if os.IsNotExist(err) {
		errMk := os.MkdirAll(dir, os.ModePerm)
		if errMk != nil {
			return "", errMk
		}
	}

	_, err = os.Stat(filePath)

	if os.IsNotExist(err) {
		jsn, errM := json.Marshal(Config{})
		if errM != nil {
			return "", errM
		}
		errW := os.WriteFile(filePath, jsn, os.ModePerm)
		if errW != nil {
			return "", errW
		}
	}

	return filePath, nil
}

func GetConfig() (Config, error) {
	cfg := Config{}

	filePath, err := getConfigFile()
	if err != nil {
		return cfg, err
	}

	dat, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(dat, &cfg); err != nil {
		return cfg, fmt.Errorf("configuration file does not have a valid format: %w", err)
	}

	return cfg, nil
}

func SetConfig(ctx context.Context, c Config) error {
	filePath, err := getConfigFile()
	if err != nil {
		return err
	}

	jsn, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsn, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
