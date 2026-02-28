package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

// BrokerMOCApp wires everything needed for commands like configuration, logging, external API client etc.
type BrokerMOCApp struct {
	APIClient *snyk.Client
	Config    *Config
	Logger    *slog.Logger
	Version   string
}

func New(cfg *Config, version string) (*BrokerMOCApp, error) {
	logger := configureLogger(cfg)

	apiClient, err := configureAPIClient(cfg, version)
	if err != nil {
		return nil, err
	}

	return &BrokerMOCApp{
		APIClient: apiClient,
		Config:    cfg,
		Logger:    logger,
		Version:   version,
	}, nil
}

func configureLogger(cfg *Config) *slog.Logger {
	var logLevel slog.Level
	switch cfg.Logging.Level {
	case "warn":
		logLevel = slog.LevelWarn
	case "info":
		logLevel = slog.LevelInfo
	case "debug":
		logLevel = slog.LevelDebug
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelWarn
	}

	handlerOpts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler
	switch cfg.Logging.Format {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	default:
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}

	return slog.New(handler)
}

func configureAPIClient(cfg *Config, version string) (*snyk.Client, error) {
	comment := "https://github.com/snyk-playground/universal-broker-multi-org-connector"
	opts := []snyk.ClientOption{snyk.WithUserAgent(fmt.Sprintf("snyk-broker-moc/%s (+%s)", version, comment))}

	if cfg.API.Region != "" {
		region, err := findRegionByAlias(cfg.API.Region)
		if err != nil {
			return nil, err
		}
		opts = append(opts, snyk.WithRegion(region))
	}

	return snyk.NewClient(cfg.API.Token, opts...)
}

func findRegionByAlias(alias string) (snyk.Region, error) {
	regions := snyk.Regions()
	// add dev as extra region
	regions = append(regions, snyk.Region{
		Alias:       "SNYK-DEV-US-01",
		AppBaseURL:  "https://app.dev.snyk.io/",
		RESTBaseURL: "https://api.dev.snyk.io/rest/",
		V1BaseURL:   "https://api.dev.snyk.io/v1/",
	})

	aliasUpper := strings.ToUpper(alias)
	for _, region := range regions {
		if strings.ToUpper(region.Alias) == aliasUpper {
			return region, nil
		}
	}

	return snyk.Region{}, fmt.Errorf("region not found: %s (available: %s)", alias, strings.Join(getRegionAliases(regions), ", "))
}

func getRegionAliases(regions []snyk.Region) []string {
	aliases := make([]string, len(regions))
	for i, r := range regions {
		aliases[i] = r.Alias
	}
	return aliases
}
