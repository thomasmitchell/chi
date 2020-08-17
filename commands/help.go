package commands

import (
	"os"
	"strings"

	"github.com/jhunt/go-ansi"
)

type helpCmd struct {
	CommandName string
}

func init() {
	Opts.Help = helpCmd{}
	Dispatch.register("help", Command{
		c:            &Opts.Help,
		usage:        "help",
		short:        "Print information about a command",
		waivedClient: true,
		waivedTarget: true,
		waivedConfig: true,
	})
}

func (c *helpCmd) ParseArgs(args []string) error {
	c.CommandName = strings.Join(args, " ")
	return nil
}

func (c *helpCmd) Run(ctx Context) error {
	if c.CommandName == "" {
		ShowGlobalHelp()
		return nil
	}

	command, err := Dispatch.Lookup(c.CommandName)
	if err != nil {
		return err
	}

	command.DisplayHelp()
	return nil
}

func ShowGlobalHelp() {
	commands := Dispatch.GetCommands()
	//Determine longest command name
	longestNameLen := 0
	for _, command := range commands {
		if len(command.name) > longestNameLen {
			longestNameLen = len(command.name)
		}
	}

	ansi.Fprintf(os.Stderr, "@R{chi} - An alternative CredHub Interface\n\n")
	//Num spaces to put between the longest command name and its description
	const minIndent = 2
	for _, command := range commands {
		numSpaces := longestNameLen - len(command.name) + minIndent
		spaceBytes := make([]byte, numSpaces)
		for i := range spaceBytes {
			spaceBytes[i] = byte(' ')
		}
		ansi.Fprintf(os.Stderr, "@C{%s}%s@M{%s}\n", command.name, spaceBytes, command.short)
	}
}
