package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

const (
	defaultConfigName     = "config"
	defaultConfigType     = "yaml"
	defaultConfigFileName = defaultConfigName + "." + defaultConfigType
	defaultEnvPrefix      = "SNYK"
)

// Config holds all configuration options needed for commands.
type Config struct {
	API     API     `mapstructure:"api"`
	Logging Logging `mapstructure:"logging"`
}

type API struct {
	Region string `mapstructure:"region"`
	Token  string `mapstructure:"token"`
} //SNYK-DEV-US-01

type Logging struct {
	Format string `mapstructure:"format"`
	Level  string `mapstructure:"level"`
}

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(brokerMOCConfigDir())
	viper.AddConfigPath(".")

	viper.SetEnvPrefix(defaultEnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			// no config at all, write default one
			err := os.MkdirAll(brokerMOCConfigDir(), 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create config folder: %w", err)
			}
			viper.SetDefault("api.region", "SNYK-US-01")
			viper.SetDefault("api.token", "")
			viper.SetDefault("logging.format", "json")
			viper.SetDefault("logging.level", "warn")

			// write config file in YAML format
			err = viper.SafeWriteConfigAs(filepath.Join(brokerMOCConfigDir(), defaultConfigFileName))
			if err != nil {
				return nil, fmt.Errorf("failed to write default config file: %w", err)
			}
		} else {
			return nil, err
		}
	}

	// workaround to bind env vars before config unmarshalling
	for _, key := range viper.AllKeys() {
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		err := viper.BindEnv(key, fmt.Sprintf("%s_%s", defaultEnvPrefix, envKey))
		if err != nil {
			return nil, err
		}
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return cfg, nil
}

func brokerMOCConfigDir() string {
	xdgConfigHomeDir := xdg.ConfigHome
	return filepath.Join(xdgConfigHomeDir, "broker-moc")
}
