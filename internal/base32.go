package internal

import (
	"encoding/base32"
	"errors"
	"flag"
	"fmt"
)

func NewBase32Command() *Base32Command {
	gc := &Base32Command{
		fs: flag.NewFlagSet("base32", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.op, "op", "encode", "encode/decode values")
	return gc
}

type Base32Command struct {
	fs *flag.FlagSet

	op string
}

func (g *Base32Command) Flag() flag.FlagSet {
	return *g.fs
}

func (g *Base32Command) Name() string {
	return g.fs.Name()
}

func (g *Base32Command) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *Base32Command) Run() (err error) {
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

func (g *Base32Command) encode(decoded string) (s string, err error) {
	s = base32.StdEncoding.EncodeToString([]byte(decoded))
	return s, err
}

func (g *Base32Command) decode(encoded string) (s string, err error) {
	data, err := base32.StdEncoding.DecodeString(encoded)
	if err != nil {
		return s, err
	}
	return string(data), err
}
