package main

import (
	"fmt"
	"os"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
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

	if commands.Opts.HelpFlag {
		args = []string{command}
		command = "help"
	}

	if command == "" {
		if len(args) > 0 {
			ansi.Fprintf(os.Stderr, "@R{unrecognized command `%s`}\n\n", args[0])
		}
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

	var client *credhub.CredHub
	if cmd.RequiresClient() {
		if !cmd.RequiresConfig() {
			panic("Command requires CredHub client but no config")
		}

		if config.Current() == nil {
			bailWith("You are not currently targeting a CredHub. Create or select one with `chi api'")
		}

		var clientTimeout = 30 * time.Second

		client, err = credhub.New(
			config.Current().Address,
			credhub.SetHttpTimeout(&clientTimeout),
			credhub.SkipTLSValidation(config.Current().SkipVerify),
			credhub.Auth(auth.Uaa("credhub_cli", "", "", "",
				config.Current().AccessToken,
				config.Current().RefreshToken,
				false),
			),
		)
		if err != nil {
			bailWith("Error creating Credhub client: %s", err)
		}
	}

	err = cmd.Run(commands.Context{
		Conf:   config,
		Client: client,
	}, args)
	if err != nil {
		bailWith(err.Error())
	}
}

func bailWith(format string, args ...interface{}) {
	ansi.Fprintf(os.Stderr, "@R{!! "+format+"}\n", args...)
	os.Exit(1)
}
