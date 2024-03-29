package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
)

func NewSha256Command() *Sha256Command {
	gc := &Sha256Command{
		fs: flag.NewFlagSet("sha256", flag.ContinueOnError),
	}

	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : <value> \n", gc.fs.Name())
	}

	return gc
}

type Sha256Command struct {
	fs *flag.FlagSet
}

func (g *Sha256Command) Flag() flag.FlagSet {
	return *g.fs
}

func (g *Sha256Command) Name() string {
	return g.fs.Name()
}

func (g *Sha256Command) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *Sha256Command) Run() (err error) {
	if g.fs.NArg() == 0 {
		return errors.New("data not provided")
	}

	s, err := g.encrypt(g.fs.Arg(0))
	if err != nil {
		return err
	}
	fmt.Println(s)
	return err
}

func (g *Sha256Command) encrypt(stringToEncrypt string) (s string, err error) {

	h := sha256.New()
	h.Write([]byte(stringToEncrypt))
	return hex.EncodeToString(h.Sum(nil)), err
}
