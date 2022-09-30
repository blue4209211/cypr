package internal

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func NewPasswordCommand() *PasswordCommand {
	gc := &PasswordCommand{
		fs: flag.NewFlagSet("password", flag.ContinueOnError),
	}

	gc.fs.IntVar(&gc.n, "n", 12, "password length")

	return gc
}

type PasswordCommand struct {
	fs *flag.FlagSet

	n int
}

func (g *PasswordCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *PasswordCommand) Name() string {
	return g.fs.Name()
}

func (g *PasswordCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *PasswordCommand) Run() (err error) {
	s, err := g.gen()
	fmt.Println(s)
	return err
}

func (g *PasswordCommand) gen() (s string, err error) {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + digits + specials

	buf := make([]byte, g.n)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < g.n; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf), err
}
