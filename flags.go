package main

import (
	"os"

	"github.com/alexflint/go-arg"
)

func parseFlags() (*AppOptions, error) {
	opts := &AppOptions{}
	p, err := arg.NewParser(arg.Config{
		Program:           AppName,
		IgnoreEnv:         true,
		IgnoreDefault:     true,
		StrictSubcommands: false,
	}, opts)
	if err != nil {
		return nil, err
	}

	p.MustParse(os.Args[1:])

	if opts.Key == "" && opts.KeyFile == "" {
		p.Fail("no --key or --keyfile flag specified")
	}

	return opts, nil
}
