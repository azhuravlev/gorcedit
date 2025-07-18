package main

const AppName = "gorcedit"

var Version = "dev" //nolint:gochecknoglobals // Defined as variable to change during build

func main() {
	app := NewApp()
	err := app.Process()
	if err != nil {
		app.Logger.Error().Err(err).Msg("Failed to edit credentials")
	}
}
