package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"lua-mountain/internal/mountain/logging"
	"lua-mountain/internal/mountain/repository"
	"lua-mountain/internal/mountain/server"
	"os"
	"path"
)

const (
	DefaultConfigPath = "/etc/mountain"
)

type (
	AppConfig struct {
		Listen server.Config `yaml:"listen"`
		Logs logging.Config `yaml:"logs"`
		Repositories []repository.Config `yaml:"repositories"`
		Storages map[string]any `yaml:"storages"`
	}

)

var (
	DefaultSearchDirs []string
	cfg               AppConfig
)

func init() {
	DefaultSearchDirs = make([]string, 0, 3)
	if wd, err := os.Getwd(); err == nil {
		DefaultSearchDirs = append(DefaultSearchDirs, wd)
	}

	if home, err := os.UserHomeDir(); err == nil {
		DefaultSearchDirs = append(DefaultSearchDirs, path.Join(home, ".mountain"))
	}

	DefaultSearchDirs = append(DefaultSearchDirs, DefaultConfigPath)
}

func Get() *AppConfig {
	return &cfg
}

func Search(dirs ...string) (string, error) {
	var (
		filenames = []string{"config.yaml", "config.yml"}
		err error
		p string
	)

	for _, dir := range dirs {
		if dir == "." {
			dir, err = os.Getwd()
			if err != nil {
				return "", err
			}
		}

		for _, filename := range filenames {
			p = path.Join(dir, filename)
			if _, err = os.Stat(p); err == nil {
				return p, nil
			}
		}
	}

	return "", fmt.Errorf("unable to search config file at dirs: %+v", dirs)
}

// Load - loads AppConfig struct from a single config file
func Load(p string) error {
	var (
		file *os.File
		err error
	)


	if file, err = os.Open(p); err != nil {
		return err
	}

	var (
		decoder = yaml.NewDecoder(file)
	)

	if err = decoder.Decode(&cfg); err != nil {
		return err
	}

	return nil
}

