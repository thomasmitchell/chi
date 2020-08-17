package commands

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/jhunt/go-ansi"
	"github.com/thomasmitchell/chi/rc"
)

type apiCmd struct {
	Name       string
	URL        string
	SkipVerify bool `cli:"-k, --skip-verify"`
}

func init() {
	Opts.API = apiCmd{}
	Dispatch.register("api", Command{
		c:            &Opts.API,
		usage:        "api [name] [url]",
		short:        "View, set, or create a CredHub target",
		waivedClient: true,
		waivedTarget: true,
	})
}

func (a *apiCmd) ParseArgs(args []string) error {
	if err := checkNumArgs(args, 0, 1, 2); err != nil {
		return err
	}
	switch len(args) {
	case 2:
		a.URL = args[1]
		fallthrough
	case 1:
		a.Name = args[0]
	}
	return nil
}

func (a *apiCmd) Run(ctx Context) error {
	if a.Name != "" {
		if a.URL != "" {
			//Create New
			err := a.createNewAPI(ctx.Conf)
			if err != nil {
				return fmt.Errorf("Error creating new API: %s", err.Error())
			}
			ctx.Conf.Save()
		} else {
			//Set Existing
			err := ctx.Conf.SetCurrent(a.Name)
			if err != nil {
				return fmt.Errorf("Error setting API: %s", err.Error())
			}
			ctx.Conf.Save()
		}
	}

	return a.displayAPI(ctx.Conf.Current())
}

func (a *apiCmd) displayAPI(target *rc.Target) error {
	if target == nil {
		return fmt.Errorf("You are not currently targeting a CredHub instance")
	}

	ansi.Printf("Currently targeting @C{%s} at @C{%s}\n", target.Name, target.Address)
	if target.SkipVerify {
		ansi.Printf("\t@R{Skipping certificate validation}\n")
	}

	return nil
}

func (a *apiCmd) createNewAPI(conf *rc.Config) error {
	err := a.canonizeURL()
	if err != nil {
		return fmt.Errorf("Error parsing input URL: %s", err.Error())
	}
	return conf.Add(rc.Target{
		Name:       a.Name,
		Address:    a.URL,
		SkipVerify: a.SkipVerify,
	})
}

var schemePrefixRegexp = regexp.MustCompile("(^[a-z]*)://")

const defaultAPIPort = "8844"

func (a *apiCmd) canonizeURL() error {
	u := a.URL
	schemeResult := schemePrefixRegexp.FindStringSubmatch(u)
	if schemeResult == nil {
		u = "https://" + u
	} else {
		scheme := string(schemeResult[1])
		if scheme != "http" && scheme != "https" {
			return fmt.Errorf("Unsupported scheme `%s'", scheme)
		}
	}

	ur, err := url.Parse(u)
	if err != nil {
		return err
	}

	if ur.Scheme == "" {
		panic("Somehow we still don't have a scheme. This is a bug.")
	}

	if ur.User != nil {
		return fmt.Errorf("URL encoded username/passwords are not supported")
	}

	if ur.Port() == "" {
		ur.Host = fmt.Sprintf("%s:%s", ur.Host, defaultAPIPort)
	}

	if ur.Path != "" && ur.Path != "/" {
		return fmt.Errorf("URL paths are not supported")
	}

	if ur.RawQuery != "" {
		return fmt.Errorf("URL parameters are not supported")
	}

	u = ur.String()
	a.URL = u
	return nil
}
