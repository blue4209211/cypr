package internal

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
)

func NewMd5Command() *Md5Command {
	gc := &Md5Command{
		fs: flag.NewFlagSet("md5", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.op, "op", "encrypt", "encrypt/decrypt values")

	return gc
}

type Md5Command struct {
	fs *flag.FlagSet

	op string
}

func (g *Md5Command) Flag() flag.FlagSet {
	return *g.fs
}

func (g *Md5Command) Name() string {
	return g.fs.Name()
}

func (g *Md5Command) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *Md5Command) Run() (err error) {
	if g.op == "encrypt" {
		if g.fs.NArg() == 0 {
			return errors.New("data not provided")
		}

		s, err := g.encrypt(g.fs.Arg(0))
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else {
		err = errors.New("Unknown Op - " + g.op)
	}
	return err
}

func (g *Md5Command) encrypt(stringToEncrypt string) (s string, err error) {

	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	data := md5.Sum(plaintext)

	return hex.EncodeToString(data[:]), err
}