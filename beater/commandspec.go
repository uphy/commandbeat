package beater

import (
	"errors"

	"github.com/uphy/commandbeat/command"
	"github.com/uphy/commandbeat/config"
	"github.com/uphy/commandbeat/parser"
)

// commandSpec implements command.Spec interface.
type commandSpec struct {
	name    string
	command string
	args    []string
	parser  parser.Parser
	debug   bool
}

func newSpec(name string, task *config.TaskConfig) (command.Spec, error) {
	c, err := task.Command()
	if err != nil {
		return nil, err
	}
	if len(c) == 0 {
		return nil, errors.New("command empty")
	}
	commandName := c[0]
	commandArgs := []string{}
	if len(c) > 0 {
		commandArgs = c[1:]
	}
	parser, err := task.Parser()
	if err != nil {
		return nil, err
	}
	return &commandSpec{name, commandName, commandArgs, parser, task.Debug}, nil
}

func (c *commandSpec) Command() string {
	return c.command
}

func (c *commandSpec) Args() []string {
	return c.args
}
