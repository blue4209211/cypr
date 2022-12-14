package internal

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
)

func NewHexCommand() *HexCommand {
	gc := &HexCommand{
		fs: flag.NewFlagSet("hex", flag.ContinueOnError),
	}

	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : <encode/decode> <value> \n", gc.fs.Name())
	}

	return gc
}

type HexCommand struct {
	fs *flag.FlagSet

	op     string
	opArgs []string
}

func (g *HexCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *HexCommand) Name() string {
	return g.fs.Name()
}

func (g *HexCommand) Init(args []string) error {
	g.op = args[0]
	err := g.fs.Parse(args)
	if err != nil {
		return err
	}
	if g.fs.NArg() != 1 {
		return errors.New("invalid args")
	}
	g.opArgs = g.fs.Args()
	return nil
}

func (g *HexCommand) Run() (err error) {
	if g.op == "encode" {
		if g.fs.NArg() == 0 {
			return errors.New("data not provided")
		}

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

func (g *HexCommand) encode(decoded string) (s string, err error) {
	s = hex.EncodeToString([]byte(decoded))
	return s, err
}

func (g *HexCommand) decode(encoded string) (s string, err error) {
	data, err := hex.DecodeString(encoded)
	if err != nil {
		return s, err
	}
	return string(data), err
}
