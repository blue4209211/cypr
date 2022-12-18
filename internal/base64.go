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

	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : <encode/decode> <value> \n", gc.fs.Name())
	}

	return gc
}

type Base64Command struct {
	fs *flag.FlagSet

	op     string
	opArgs []string
}

func (g *Base64Command) Flag() flag.FlagSet {
	return *g.fs
}

func (g *Base64Command) Name() string {
	return g.fs.Name()
}

func (g *Base64Command) Init(args []string) error {
	g.op = args[0]
	err := g.fs.Parse(args[1:])
	if err != nil {
		return err
	}
	if g.fs.NArg() != 1 {
		return errors.New("invalid args")
	}
	g.opArgs = g.fs.Args()
	return nil
}

func (g *Base64Command) Run() (err error) {
	if g.op == "encode" {
		s, err := g.encode(g.opArgs[0])
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else if g.op == "decode" {
		s, err := g.decode(g.opArgs[0])
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
	s = base64.StdEncoding.EncodeToString([]byte(decoded))
	return s, err
}

func (g *Base64Command) decode(encoded string) (s string, err error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return s, err
	}
	return string(data), err
}
