package internal

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func NewRandCommand() *RandCommand {
	gc := &RandCommand{
		fs: flag.NewFlagSet("rand", flag.ContinueOnError),
	}

	gc.fs.IntVar(&gc.n, "n", 12, "rand length")
	gc.fs.IntVar(&gc.a, "a", -1, "alphabet length")
	gc.fs.IntVar(&gc.d, "d", -1, "digit length")
	gc.fs.IntVar(&gc.s, "s", -1, "special char length")

	return gc
}

type RandCommand struct {
	fs *flag.FlagSet

	n int
	a int
	d int
	s int
}

func (g *RandCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *RandCommand) Name() string {
	return g.fs.Name()
}

func (g *RandCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *RandCommand) Run() (err error) {
	s, err := g.gen()
	fmt.Println(s)
	return err
}

func (g *RandCommand) gen() (s string, err error) {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	alphabets := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz"
	all := ""
	if g.a != 0 {
		all = all + alphabets
	}
	if g.n != 0 {
		all = all + digits
	}
	if g.s != 0 {
		all = all + specials
	}

	v := 0
	if g.a > 0 {
		v = g.a
		if g.a == g.n {
			g.s = 0
			g.d = 0
		}
	}
	if g.d > 0 {
		v = g.d
		if g.d == g.n {
			g.a = 0
			g.s = 0
		}
	}
	if g.s > 0 {
		v = g.s
		if g.s == g.n {
			g.a = 0
			g.d = 0
		}
	}

	if v > g.n {
		return s, errors.New("overall sum is greater than n")
	}

	buf := make([]byte, g.n)
	c := 0
	if g.d != 0 {
		if g.d > 0 {
			for i := c; i < g.d; i++ {
				buf[i] = digits[rand.Intn(len(digits))]
				c = c + 1
			}
		} else {
			buf[0] = digits[rand.Intn(len(digits))]
			c = c + 1
		}
	}
	if g.s != 0 {
		if g.s > 0 {
			for i := c; i < g.s; i++ {
				buf[i] = specials[rand.Intn(len(specials))]
				c = c + 1
			}
		} else {
			buf[c] = specials[rand.Intn(len(specials))]
			c = c + 1
		}
	}
	if g.a != 0 {
		if g.a > 0 {
			for i := c; i < g.a; i++ {
				buf[i] = alphabets[rand.Intn(len(alphabets))]
				c = c + 1
			}
		} else {
			buf[c] = alphabets[rand.Intn(len(alphabets))]
			c = c + 1
		}

	}

	for i := c; i < g.n; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}

	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})

	return string(buf), err
}
