package commands

import (
	"fmt"
	"os"

	"github.com/thomasmitchell/chi/commands/internal"
	"gopkg.in/yaml.v2"
)

type getCmd struct {
	Path    string
	Verbose bool `cli:"-v, --verbose"`
}

func init() {
	Opts.Get = getCmd{}
	Dispatch.register("get", Command{
		c:     &Opts.Get,
		usage: "get <path>",
		short: "Get the value of a secret",
	})
}

func (p *getCmd) ParseArgs(args []string) error {
	if err := checkNumArgs(args, 1); err != nil {
		return err
	}

	p.Path = "/"
	if len(args) > 0 {
		p.Path = internal.CanonizePathForAPI(args[0])
	}

	return nil
}

func (p *getCmd) Run(ctx Context) error {
	cred, err := ctx.Client.GetLatestVersion(p.Path)
	if err != nil {
		return fmt.Errorf("Error fetching secret: %s", err.Error())
	}

	var toEncode interface{} = cred

	if !p.Verbose {
		toEncode = cred.Value
	}
	fmt.Printf("--- # %s\n", internal.CanonizePathForOutput(p.Path))
	enc := yaml.NewEncoder(os.Stdout)
	err = enc.Encode(toEncode)
	return err
}
