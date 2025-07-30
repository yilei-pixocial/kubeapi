package sysinit

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/olebedev/config"
)

var GCF *config.Config //global config

func InitConf() error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	configPath := filepath.Join(pwd, "configs", "application.yml")
	cfg, err := config.ParseYamlFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	GCF = cfg
	return nil
}
