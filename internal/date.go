package internal

import (
	"errors"
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

func parseDateOrTime(s string) (t time.Time, f string, err error) {
	formats := map[string]string{
		"2006-01-02T15:04:05Z07:00": "datetime",
		"2006-01-02T15:04:05 MST":   "datetime",
		"3:04PM":                    "time",
		"3:04PM MST":                "time",
		"15:04 MST":                 "time",
		"15:04":                     "time",
	}
	for k, f := range formats {
		if f == "time" {
			t, err = time.Parse("2006-01-02 "+k, "2006-01-02 "+s)
		} else {
			t, err = time.Parse(k, s)
		}
		if err == nil {
			return t, f, err
		}
	}
	return t, f, errors.New("unable to parse date/time")
}

func (g *DateCommand) Run() (err error) {
	current := time.Now()
	currentTz, _ := current.Zone()
	format := "datetime"
	if len(g.opArgs) > 0 {
		nixTime, err := strconv.ParseInt(g.opArgs[0], 10, 64)
		if err == nil {
			current = time.UnixMilli(nixTime)
		} else {
			current2, f, err := parseDateOrTime(g.opArgs[0])
			if err != nil {
				return err
			}
			current = current2.In(current.Location())
			format = f
		}
	}

	timeFormat := "15:04"
	if format == "datetime" {
		fmt.Printf("Local(%v) - %v \n", currentTz, current)
		fmt.Printf("UTC - %v \n", current.UTC())
		fmt.Printf("Unix Milli - %v \n", current.UnixMilli())
	} else {
		fmt.Printf("Local(%v) - %v \n", currentTz, current.Format(timeFormat))
		fmt.Printf("UTC - %v \n", current.UTC().Format(timeFormat))
	}
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
			if format == "datetime" {
				fmt.Printf("%s - %v \n", t, c)
			} else {
				fmt.Printf("%s - %v \n", t, c.Format(timeFormat))
			}
		}
	}

	return err
}
