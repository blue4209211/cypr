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
	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : <encode/decode> <value> \n", gc.fs.Name())
	}
	return gc
}

type Base32Command struct {
	fs *flag.FlagSet

	op     string
	opArgs []string
}

func (g *Base32Command) Flag() flag.FlagSet {
	return *g.fs
}

func (g *Base32Command) Name() string {
	return g.fs.Name()
}

func (g *Base32Command) Init(args []string) error {
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

func (g *Base32Command) Run() (err error) {
	switch g.op {
	case "encode":
		s, err := g.encode(g.opArgs[0])
		if err != nil {
			return err
		}
		fmt.Println(s)
	case "decode":
		s, err := g.decode(g.opArgs[0])
		if err != nil {
			return err
		}
		fmt.Println(s)
	default:
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
