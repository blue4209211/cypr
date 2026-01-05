package internal

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
)

func NewUrlCommand() *UrlCommand {
	gc := &UrlCommand{
		fs: flag.NewFlagSet("url", flag.ContinueOnError),
	}

	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : <encode/decode> <value> \n", gc.fs.Name())
	}

	return gc
}

type UrlCommand struct {
	fs *flag.FlagSet

	op   string
	args []string
}

func (g *UrlCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *UrlCommand) Name() string {
	return g.fs.Name()
}

func (g *UrlCommand) Init(args []string) error {
	g.op = args[0]
	err := g.fs.Parse(args[1:])
	if err != nil {
		return err
	}
	if g.fs.NArg() != 1 {
		return errors.New("invalid args")
	}

	g.args = g.fs.Args()
	return nil
}

func (g *UrlCommand) Run() (err error) {
	switch g.op {
	case "encode":
		s, err := g.encode(g.args[0])
		if err != nil {
			return err
		}
		fmt.Println(s)
	case "decode":
		s, err := g.decode(g.args[0])
		if err != nil {
			return err
		}
		fmt.Println(s)
	default:
		err = errors.New("Unknown Op - " + g.op)
	}
	return err
}

func (g *UrlCommand) encode(rawString string) (s string, err error) {
	return url.QueryEscape(rawString), err
}

func (g *UrlCommand) decode(escaped string) (s string, err error) {
	return url.QueryUnescape(escaped)
}
