package main

import (
	"bytes"
	"fmt"
	"os"
)

// Process read, decode, edit, encode and save the credentials
func (app *App) Process() error {
	app.Logger.Info().Msg("start editing")

	tmpFile, err := os.CreateTemp("", "*.credentials")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	defer func() {
		tmpErr := os.Remove(tmpFile.Name())
		if tmpErr != nil {
			app.Logger.Error().Err(tmpErr).Msg("Failed to remove temporary file")
		}
	}()

	decoded, err := app.decodeFile()
	if err != nil {
		return fmt.Errorf("failed to decode credentials: %w", err)
	}

	if _, err = tmpFile.Write(decoded); err != nil {
		return fmt.Errorf("failed to write decrypted credentials: %w", err)
	}

	err = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	err = EditFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to edit credentials with external editor: %w", err)
	}

	updated, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read updated credentials: %w", err)
	}
	app.Logger.Debug().Bytes("data", updated).Msg("edited credentials")

	if bytes.Equal(decoded, updated) {
		app.Logger.Info().Msg("No changes were made, credentials file was not updated")
		return nil
	}

	err = app.encodeFile(updated)
	if err != nil {
		return fmt.Errorf("failed to write updated credentials: %w", err)
	}

	app.Logger.Info().Msg("edited credentials file encrypted and saved")
	return nil
}

func (app *App) decodeFile() ([]byte, error) {
	encrypted, err := os.ReadFile(app.opts.CredsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open credentials: %w", err)
	}
	app.Logger.Debug().Bytes("data", encrypted).Msg("encrypted credentials")

	var decrypted []byte
	if len(encrypted) > 0 {
		decrypted, err = Decrypt(string(encrypted), app.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt credentials: %w", err)
		}
	}
	app.Logger.Debug().Bytes("data", decrypted).Msg("decrypted credentials")

	decoded, err := UnmarshalRubyString(decrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted credentials: %w", err)
	}
	app.Logger.Debug().Bytes("data", decoded).Msg("decoded credentials")

	return decoded, nil
}

func (app *App) encodeFile(data []byte) error {
	strToSave, err := MarshalRubyString(data)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}
	app.Logger.Debug().Bytes("data", strToSave).Msg("encoded credentials")

	credsToSave, err := Encrypt(strToSave, app.Key)
	if err != nil {
		return fmt.Errorf("failed to encrypt credentials: %w", err)
	}

	if err = os.Truncate(app.opts.CredsPath, 0); err != nil {
		return fmt.Errorf("failed to truncate credentials file: %w", err)
	}

	if err = os.WriteFile(app.opts.CredsPath, []byte(credsToSave), app.encPerms); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}
