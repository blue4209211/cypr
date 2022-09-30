package internal

import (
	"flag"
	"fmt"

	"github.com/google/uuid"
)

func NewUuidCommand() *UuidCommand {
	gc := &UuidCommand{
		fs: flag.NewFlagSet("uuid", flag.ContinueOnError),
	}

	return gc
}

type UuidCommand struct {
	fs *flag.FlagSet
}

func (g *UuidCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *UuidCommand) Name() string {
	return g.fs.Name()
}

func (g *UuidCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *UuidCommand) Run() (err error) {
	s := uuid.New()
	fmt.Println(s.String())
	return err
}
