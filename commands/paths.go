package commands

import (
	"fmt"
	"sort"

	"github.com/thomasmitchell/chi/commands/internal"
)

type pathsCmd struct {
	Path string
}

func init() {
	Opts.Paths = pathsCmd{}
	Dispatch.register("paths", Command{
		c:     &Opts.Paths,
		usage: "paths [path]",
		short: "List secret paths",
	})
}

func (p *pathsCmd) ParseArgs(args []string) error {
	if err := checkNumArgs(args, 0, 1); err != nil {
		return err
	}

	p.Path = "/"
	if len(args) > 0 {
		p.Path = internal.CanonizePathForAPI(args[0])
	}

	return nil
}

func (p *pathsCmd) Run(ctx Context) error {
	results, err := ctx.Client.FindByPath(p.Path)
	if err != nil {
		return fmt.Errorf("Error fetching paths: %s", err.Error())
	}

	sort.Slice(results.Credentials, func(i, j int) bool {
		return results.Credentials[i].Name < results.Credentials[j].Name
	})

	for _, cred := range results.Credentials {
		fmt.Println(internal.CanonizePathForOutput(cred.Name))
	}

	return nil
}
