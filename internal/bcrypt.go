package internal

import (
	"errors"
	"flag"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func NewBcryptCommand() *BcryptCommand {
	gc := &BcryptCommand{
		fs: flag.NewFlagSet("bcrypt", flag.ContinueOnError),
	}

	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : <hash/verify> [args] \n", gc.fs.Name())
		fmt.Printf(" hash: <password>\n")
		fmt.Printf(" verify: <hash> <password>\n")
	}

	return gc
}

type BcryptCommand struct {
	fs     *flag.FlagSet
	op     string
	opArgs []string
}

func (g *BcryptCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *BcryptCommand) Name() string {
	return g.fs.Name()
}

func (g *BcryptCommand) Init(args []string) error {
	if len(args) < 1 {
		return errors.New("invalid args, op required")
	}
	g.op = args[0]

	err := g.fs.Parse(args[1:])
	if err != nil {
		return err
	}

	if g.op == "hash" {
		if g.fs.NArg() != 1 {
			return errors.New("invalid args, password required")
		}
	} else if g.op == "verify" {
		if g.fs.NArg() != 2 {
			return errors.New("invalid args, hash and password required")
		}
	} else {
		return fmt.Errorf("unknown operation: %s", g.op)
	}

	g.opArgs = g.fs.Args()
	return nil
}

func (g *BcryptCommand) Run() error {
	if g.op == "hash" {
		return g.hash(g.opArgs[0])
	} else if g.op == "verify" {
		return g.verify(g.opArgs[0], g.opArgs[1])
	}
	return nil
}

func (g *BcryptCommand) hash(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	fmt.Println(string(hash))
	return nil
}

func (g *BcryptCommand) verify(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Invalid password")
		return err
	}
	fmt.Println("Valid password")
	return nil
}
