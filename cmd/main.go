package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/blue420211/cypr/internal"
)

func root(args []string) (cmd internal.Command, e error) {
	cmds := []internal.Command{
		internal.NewAesCommand(),
		internal.NewBase32Command(),
		internal.NewBase64Command(),
		internal.NewHexCommand(),
		internal.NewMd5Command(),
		internal.NewSha1Command(),
		internal.NewSha256Command(),
		internal.NewRandCommand(),
		internal.NewUuidCommand(),
		internal.NewUrlCommand(),
		internal.NewDateCommand(),
		internal.NewJwtCommand(),
	}

	cmdsStr := []string{}
	for _, c := range cmds {
		cmdsStr = append(cmdsStr, c.Name())
	}
	cmdsStr = append(cmdsStr, "help")

	if len(args) < 1 {
		return cmd, errors.New("you must pass a sub-command - " + strings.Join(cmdsStr, ","))
	}

	subcommand := os.Args[1]

	if subcommand == "help" {
		return cmd, errors.New("you must pass a sub-command - " + strings.Join(cmdsStr, ","))
	}

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			err := cmd.Init(os.Args[2:])
			if err != nil {
				return cmd, fmt.Errorf("%v", err)
			}
			return cmd, cmd.Run()
		}
	}

	return cmd, fmt.Errorf("unknown subcommand: %s", subcommand)
}

func main() {

	if c, err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		fmt.Println("-------------------")
		if c != nil {
			c.Flag().Usage()
		}
		os.Exit(-1)
	}
}
