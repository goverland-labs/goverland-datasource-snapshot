package cli

import (
	"github.com/akamensky/argparse"
)

type ArgumentsDetails map[string]string
type Arguments map[string]*string

func (a *Arguments) Get(name string) string {
	v, ok := (*a)[name]
	if !ok {
		return ""
	}

	return *v
}

type Command interface {
	ParseArgs(args ...string) (Arguments, error)
	GetName() string
	Execute(args Arguments) error
	GetArguments() ArgumentsDetails
}

type base struct {
}

func (base) parseArgs(c Command, args ...string) (map[string]*string, error) {
	parser := argparse.NewParser(c.GetName(), "")

	values := make(map[string]*string)

	for name := range c.GetArguments() {
		values[name] = parser.String("", name, nil)
	}

	if err := parser.Parse(append([]string{c.GetName()}, args...)); err != nil {
		return nil, err
	}

	return values, nil
}
