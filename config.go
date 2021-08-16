package appfx

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/config"
	"go.uber.org/fx"
)

var configOptions = fx.Options(
	fx.Provide(provideConfig),
)

func provideConfig() (config.Provider, error) {
	yamlFilesOpts := func(dir string) []config.YAMLOption {
		files := []config.YAMLOption{}

		des, err := os.ReadDir(dir)
		if err != nil {
			return nil
		}

		for _, de := range des {
			if !de.IsDir() && (strings.HasSuffix(de.Name(), ".yml") || strings.HasSuffix(de.Name(), ".yaml")) {
				files = append(files, config.File(fmt.Sprintf("./%s", de.Name())))
			}
		}

		return files
	}

	d := os.Getenv("APP_ADDITIONAL_CONFIGS_DIR")
	confdir := os.Getenv("APP_CONFIGS_DIR")
	if confdir == "" {
		confdir = "/etc/app"
	}

	fileOpts := yamlFilesOpts(confdir)
	fileOpts = append(fileOpts, yamlFilesOpts(".")...)
	fileOpts = append(fileOpts, yamlFilesOpts(d)...)

	return config.NewYAML(fileOpts...)
}
