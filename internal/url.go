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

	gc.fs.StringVar(&gc.op, "op", "encode", "encrypt/decrypt values")

	return gc
}

type UrlCommand struct {
	fs *flag.FlagSet

	op string
}

func (g *UrlCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *UrlCommand) Name() string {
	return g.fs.Name()
}

func (g *UrlCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *UrlCommand) Run() (err error) {
	if g.op == "encode" {
		if g.fs.NArg() == 0 {
			return errors.New("data not provided")
		}

		s, err := g.encode(g.fs.Arg(0))
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else if g.op == "decode" {
		s, err := g.decode(g.fs.Arg(0))
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else {
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
