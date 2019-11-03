package main

import (
	"fmt"
	"os"

	"github.com/jhunt/go-ansi"
	"github.com/jhunt/go-cli"
	"github.com/thomasmitchell/chi/commands"
	"github.com/thomasmitchell/chi/rc"
)

var cliConf *rc.Config

func main() {
	command, args, err := cli.Parse(&commands.Opts)
	if err != nil {
		bailWith(err.Error())
	}

	if commands.Opts.VersionFlag {
		args = nil
		command = "version"
	} else if commands.Opts.HelpFlag {
		args = []string{command}
		command = "help"
	}

	if command == "" {
		commands.ShowGlobalHelp()
		os.Exit(1)
	}

	cmd, err := commands.Dispatch.Lookup(command)
	if err != nil {
		commands.ShowGlobalHelp()
		fmt.Fprintf(os.Stderr, "\n")
		bailWith(err.Error())
	}

	var config *rc.Config
	if cmd.RequiresConfig() {
		config, err = rc.Load(rc.DefaultPath)
		if err != nil {
			bailWith("Error loading config: %s", err.Error())
		}
	}

	err = cmd.Run(commands.Context{
		Conf: config,
	}, args)
	if err != nil {
		bailWith(err.Error())
	}
}

func bailWith(format string, args ...interface{}) {
	ansi.Fprintf(os.Stderr, "@R{!! "+format+"}", args...)
	os.Exit(1)
}
