package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// App contains everything required for the execution
type App struct {
	Logger   zerolog.Logger
	Key      []byte
	encPerms os.FileMode
	opts     *AppOptions
}

// NewApp builds and configures the App
func NewApp() *App {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(output).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	opts, err := parseFlags()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not parse cmd flags")
	}
	if opts.Debug {
		logger = logger.Level(zerolog.DebugLevel)
	}

	credsInfo, exists, err := FileInfo(opts.CredsPath)
	if !exists {
		logger.Fatal().Err(err).Msg("could not find credentials file")
	}

	if err != nil {
		logger.Fatal().Err(err).Msg("could not read credentials file flags")
	}

	app := &App{Logger: logger, encPerms: credsInfo.Mode(), opts: opts}
	app.decodeEncKey()
	return app
}

func (app *App) decodeEncKey() {
	keyStr, err := app.keyString()
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("can't read key")
	}

	key, err := hex.DecodeString(keyStr)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("invalid key passed")
	}
	app.Key = key
}

func (app *App) keyString() (string, error) {
	if app.opts.Key != "" {
		return strings.TrimSpace(app.opts.Key), nil
	}
	keyBytes, err := os.ReadFile(app.opts.KeyFile)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %w", err)
	}
	return strings.TrimSpace(string(keyBytes)), nil
}
