package main

import "fmt"

// AppOptions configuration definition
type AppOptions struct {
	CredsPath string `arg:"positional" help:"Path to encrypted credentials file"`
	Key       string `arg:"-k,--key"  help:"Key value"`
	KeyFile   string `arg:"-f,--keyfile"  help:"Path to file contains key"`
	Debug     bool   `arg:"-d"  help:"Turn on debug logging"`
}

// Version returns version for help
func (opt AppOptions) Version() string {
	return fmt.Sprintf("%s (%s)", AppName, Version)
}

// Description contains app description for help
func (opt AppOptions) Description() string {
	return "Simple Rails encrypted credentials editor, works only with strings encoded by Ruby Marshal 4.8."
}
