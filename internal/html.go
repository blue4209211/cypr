package internal

import (
	"errors"
	"flag"
	"fmt"
	"html"
)

func NewHtmlCommand() *HtmlCommand {
	gc := &HtmlCommand{
		fs: flag.NewFlagSet("html", flag.ContinueOnError),
	}

	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : <encode/decode> <value> \n", gc.fs.Name())
	}

	return gc
}

type HtmlCommand struct {
	fs *flag.FlagSet

	op   string
	args []string
}

func (g *HtmlCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *HtmlCommand) Name() string {
	return g.fs.Name()
}

func (g *HtmlCommand) Init(args []string) error {
	if len(args) == 0 {
		return errors.New("missing operation")
	}
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

func (g *HtmlCommand) Run() (err error) {
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

func (g *HtmlCommand) encode(rawString string) (s string, err error) {
	return html.EscapeString(rawString), nil
}

func (g *HtmlCommand) decode(escaped string) (s string, err error) {
	return html.UnescapeString(escaped), nil
}
