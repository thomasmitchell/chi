package commands

import (
	"fmt"
	"os"

	"github.com/doomsday-project/doomsday/storage/uaa"
	"github.com/jhunt/go-ansi"
	"github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh/terminal"
)

type loginCmd struct {
	Username string `cli:"-u, --username"`
	Password string `cli:"-p, --password"`
	//Currently only supports implicit auth with the default client id/secret
}

func init() {
	Opts.Login = loginCmd{}
	Dispatch.register("login", Command{
		c:     &Opts.Login,
		usage: "login [-u <username>] [-p <password>]",
		short: "Log-in to the targeted CredHub API",
	})
}

func (l *loginCmd) ParseArgs(args []string) error {
	return checkNumArgs(args, 0)
}

func (l *loginCmd) Run(ctx Context) error {
	const (
		DefaultClientID     = "credhub_cli"
		DefaultClientSecret = ""
	)

	authURL, err := ctx.Client.AuthURL()
	if err != nil {
		return fmt.Errorf("Error getting auth URL: %s", err)
	}

	currentCfg := ctx.Conf.Current()

	uaaClient := uaa.Client{
		URL:               authURL,
		SkipTLSValidation: currentCfg.SkipVerify,
	}

	isTerminal := isatty.IsTerminal(os.Stdin.Fd())
	username := l.Username
	if username == "" && isTerminal {
		username, err = l.promptUsername()
		if err != nil {
			return fmt.Errorf("Error prompting for username: %s", err)
		}
	}

	password := l.Password
	if password == "" && isTerminal {
		password, err = l.promptPassword()
		if err != nil {
			return fmt.Errorf("Error prompting for password: %s", err)
		}
	}

	authResp, err := uaaClient.Password(
		DefaultClientID,
		DefaultClientSecret,
		username,
		password,
	)
	if err != nil {
		return err
	}

	currentCfg.AccessToken = authResp.AccessToken
	currentCfg.RefreshToken = authResp.RefreshToken
	err = ctx.Conf.Save()
	if err != nil {
		return fmt.Errorf("Error saving authentication: %s", err)
	}

	return nil
}

func (l *loginCmd) promptUsername() (string, error) {
	var username string
	fmt.Fprintf(os.Stderr, "Username: ")
	_, err := fmt.Scanln(&username)
	return username, err
}

func (l *loginCmd) promptPassword() (string, error) {
	ansi.Fprintf(os.Stderr, "Password: ")
	passwordBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	ansi.Fprintf(os.Stderr, "\n")
	return string(passwordBytes), nil
}
