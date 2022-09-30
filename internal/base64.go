package internal

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
)

func NewBase64Command() *Base64Command {
	gc := &Base64Command{
		fs: flag.NewFlagSet("base64", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.op, "op", "encode", "encode/decode values")
	return gc
}

type Base64Command struct {
	fs *flag.FlagSet

	op string
}

func (g *Base64Command) Flag() flag.FlagSet {
	return *g.fs
}

func (g *Base64Command) Name() string {
	return g.fs.Name()
}

func (g *Base64Command) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *Base64Command) Run() (err error) {
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

func (g *Base64Command) encode(decoded string) (s string, err error) {
	s = base64.RawStdEncoding.EncodeToString([]byte(decoded))
	return s, err
}

func (g *Base64Command) decode(encoded string) (s string, err error) {
	data, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return s, err
	}
	return string(data), err
}
