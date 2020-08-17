package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"github.com/jhunt/go-ansi"
	"github.com/thomasmitchell/chi/rc"
)

type Context struct {
	//Conf should be set to the rc.Config that was loaded from the configuration file.
	Conf   *rc.Config
	Client *credhub.CredHub
}

//FLAGS
type Options struct {
	HelpFlag    bool `cli:"-h, --help"`
	VersionFlag bool `cli:"-v, --version"`

	API     apiCmd     `cli:"api"`
	Login   loginCmd   `cli:"login"`
	Paths   pathsCmd   `cli:"paths"`
	Version versionCmd `cli:"version"`
	Help    helpCmd    `cli:"help"`
}

var Opts = Options{}

//COMMAND
type Command struct {
	c            cmd
	name         string
	usage        string
	short        string
	long         string
	waivedClient bool
	waivedTarget bool
	waivedConfig bool
}

func (c *Command) Run(ctx Context, args []string) error {
	if err := c.c.ParseArgs(args); err != nil {
		return err
	}
	return c.c.Run(ctx)
}

func (c Command) RequiresClient() bool {
	return !c.waivedClient
}

func (c Command) RequiresConfig() bool {
	return !c.waivedConfig
}

func (c Command) DisplayHelp() {
	ansi.Fprintf(os.Stderr, "@C{%s}\n\n@G{USAGE:} @M{%s}\n\n%s\n", c.short, c.usage, c.long)
}

type cmd interface {
	ParseArgs([]string) error
	Run(Context) error
}

//DISPATCHER
type Dispatcher struct {
	cmdMap  map[string]*Command
	aliases map[string]*Command
}

var Dispatch = &Dispatcher{
	cmdMap: map[string]*Command{},
}

func (d *Dispatcher) register(name string, c Command) {
	c.name = name
	d.cmdMap[name] = &c
}

func (d *Dispatcher) alias(name, aliasFor string) {
	cmd, err := d.Lookup(aliasFor)
	if err != nil {
		panic(fmt.Sprintf("attempting to alias for unknown command: `%s'", aliasFor))
	}

	d.aliases[name] = cmd
}

func (d *Dispatcher) Lookup(name string) (*Command, error) {
	ret := d.cmdMap[name]
	if ret == nil {
		ret = d.aliases[name]
		if ret == nil {
			return nil, fmt.Errorf("Unknown command: %s", name)
		}
	}
	return ret, nil
}

//GetCommands returns a slice of commands, sorted lexicographically by the name
// of each command.
func (d *Dispatcher) GetCommands() []*Command {
	commands := make([]*Command, 0, len(d.cmdMap))
	for _, c := range d.cmdMap {
		commands = append(commands, c)
	}

	sort.Slice(commands, func(i, j int) bool { return commands[i].name < commands[j].name })
	return commands
}

//HELPERS
func checkNumArgs(args []string, allowed int, moreAllowed ...int) error {
	found := false
	numArgs := len(args)
	for _, n := range append(moreAllowed, allowed) {
		if numArgs == n {
			found = true
			break
		}
	}

	var err error
	if !found {
		err = formatNumArgsError(numArgs, allowed, moreAllowed...)
	}
	return err
}

func formatNumArgsError(numArgs int, allowed int, moreAllowed ...int) error {
	sort.Ints(moreAllowed)
	possibilityStr := ""
	if len(moreAllowed) == 1 {
		possibilityStr = fmt.Sprintf(" or %d", moreAllowed)
	}

	if len(moreAllowed) > 1 {
		moreAllowedStrings := make([]string, 0, len(moreAllowed)-1)
		for _, a := range moreAllowed[:len(moreAllowed)-1] {
			moreAllowedStrings = append(moreAllowedStrings, fmt.Sprintf("%d", a))
		}

		possibilityStr = fmt.Sprintf(", %s, or %d", strings.Join(moreAllowedStrings, ", "), moreAllowed[len(moreAllowed)-1])
	}

	errorMessage := fmt.Sprintf("Expected %d%s args: got %d", allowed, possibilityStr, numArgs)
	return fmt.Errorf(errorMessage)
}
