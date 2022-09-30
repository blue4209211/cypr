package internal

import "flag"

type Command interface {
	Init([]string) error
	Run() error
	Name() string
	Flag() flag.FlagSet
}
