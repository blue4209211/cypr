package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/blue420211/cypr/internal"
)

func root(args []string) (cmd internal.Command, e error) {
	cmds := []internal.Command{
		internal.NewAesCommand(),
		internal.NewBase32Command(),
		internal.NewBase64Command(),
		internal.NewHexCommand(),
		internal.NewMd5Command(),
		internal.NewPasswordCommand(),
		internal.NewUuidCommand(),
	}

	if len(args) < 1 {
		return cmd, errors.New("you must pass a sub-command - ")
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd, cmd.Run()
		}
	}

	return cmd, fmt.Errorf("unknown subcommand: %s", subcommand)
}

func main() {

	if c, err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		fmt.Println("-------------------")
		c.Flag().Usage()
	}
}
