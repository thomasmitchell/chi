package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/thomasmitchell/chi/commands/internal"
	"gopkg.in/yaml.v2"
)

type getCmd struct {
	Path    string
	Subpath internal.Subpath
	Verbose bool `cli:"-v, --verbose"`
}

func init() {
	Opts.Get = getCmd{}
	Dispatch.register("get", Command{
		c:     &Opts.Get,
		usage: "get <path>[:<subpath>]",
		short: "Get the value of a secret",
	})
}

func (p *getCmd) ParseArgs(args []string) error {
	if err := checkNumArgs(args, 1); err != nil {
		return err
	}

	path, subpath := internal.SplitPath(args[0])
	p.Path = internal.CanonizePathForAPI(path)
	p.Subpath = internal.NewSubpath(subpath)
	if !p.Verbose {
		p.Subpath = internal.NewSubpath(strings.Join([]string{"value", string(p.Subpath)}, "."))
	}

	return nil
}

func (p *getCmd) Run(ctx Context) error {
	cred, err := ctx.Client.GetLatestVersion(p.Path)
	if err != nil {
		return fmt.Errorf("Error fetching secret: %s", err.Error())
	}

	secret, err := internal.NewSecretFromCredential(cred)
	if err != nil {
		return fmt.Errorf("Error parsing secret: %s", err.Error())
	}

	toEncode, err := secret.GetSubpath(p.Subpath)
	if err != nil {
		return fmt.Errorf("Error filtering to subpath: %s", err.Error())
	}

	switch v := toEncode.(type) {
	case string, float64, bool:
		fmt.Printf("%+v\n", v)
	default:
		enc := yaml.NewEncoder(os.Stdout)
		err = enc.Encode(toEncode)
	}

	return err
}
