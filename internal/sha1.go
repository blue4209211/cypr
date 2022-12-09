package internal

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
)

func NewSha1Command() *Sha1Command {
	gc := &Sha1Command{
		fs: flag.NewFlagSet("sha1", flag.ContinueOnError),
	}

	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : <value> \n", gc.fs.Name())
	}

	return gc
}

type Sha1Command struct {
	fs *flag.FlagSet
}

func (g *Sha1Command) Flag() flag.FlagSet {
	return *g.fs
}

func (g *Sha1Command) Name() string {
	return g.fs.Name()
}

func (g *Sha1Command) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *Sha1Command) Run() (err error) {
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

func (g *Sha1Command) encrypt(stringToEncrypt string) (s string, err error) {

	plaintext := []byte(stringToEncrypt)

	data := sha1.Sum(plaintext)

	return hex.EncodeToString(data[:]), err
}
