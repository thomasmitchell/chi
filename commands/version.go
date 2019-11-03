package commands

import (
	"fmt"

	"github.com/thomasmitchell/chi/version"
)

type versionCmd struct{}

func init() {
	Opts.Version = versionCmd{}
	Dispatch.register("version", Command{
		c:            &Opts.Version,
		usage:        "version",
		short:        "Get the version of this program",
		waivedClient: true,
		waivedTarget: true,
		waivedConfig: true,
	})
}

func (*versionCmd) ParseArgs(args []string) error { return checkNumArgs(args, 0) }

func (*versionCmd) Run(ctx Context) error {
	fmt.Printf("chi version %s\n", version.Version)
	return nil
}
