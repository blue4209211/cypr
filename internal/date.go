package internal

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tkuchiki/go-timezone"
)

func NewDateCommand() *DateCommand {
	gc := &DateCommand{
		fs: flag.NewFlagSet("date", flag.ContinueOnError),
	}

	gc.fs.Usage = func() {
		fmt.Printf("Usage of %s : [args] [value] \n Args:\n", gc.fs.Name())
		gc.fs.PrintDefaults()
	}

	gc.fs.StringVar(&gc.tz, "tz", "", "Destination Timezones, optional")
	// gc.fs.StringVar(&gc.nonce, "nonce", "", "12 bytes nonce value (optional) if not provided then random nonce will be used and appended to cypher text")

	return gc
}

type DateCommand struct {
	fs *flag.FlagSet

	tz     string
	opArgs []string
}

func (g *DateCommand) Flag() flag.FlagSet {
	return *g.fs
}

func (g *DateCommand) Name() string {
	return g.fs.Name()
}

func (g *DateCommand) Init(args []string) error {
	err := g.fs.Parse(args)
	if err != nil {
		return err
	}
	g.opArgs = g.fs.Args()
	return nil
}

func (g *DateCommand) Run() (err error) {
	current := time.Now()
	if len(g.opArgs) > 0 {
		nixTime, err := strconv.ParseInt(g.opArgs[0], 10, 64)
		if err == nil {
			current = time.UnixMilli(nixTime)
		} else {
			current, err = time.Parse(time.RFC3339, g.opArgs[0])
			if err != nil {
				return err
			}

		}
	}

	currentTz, _ := current.Zone()
	fmt.Printf("Local(%v) - %v \n", currentTz, current)
	fmt.Printf("UTC - %v \n", current.UTC())
	fmt.Printf("Unix Milli - %v \n", current.UnixMilli())
	if len(g.tz) > 0 {
		fmt.Printf("---------- \n")
		for _, t := range strings.Split(g.tz, ",") {
			tzs := t
			tz := timezone.New()
			if len(t) <= 3 {
				if err != nil {
					return err
				}

				tzInfo, err := tz.GetTimezones(t)
				if err != nil {
					return err
				}
				tzs = tzInfo[0]
			}
			c, err := tz.FixedTimezone(current, tzs)
			if err != nil {
				return err
			}
			fmt.Printf("%s - %v \n", t, c)
		}
	}

	return err
}
