package beater

import (
	"errors"
	"fmt"

	"github.com/uphy/commandbeat/command"
	"github.com/uphy/commandbeat/config"
	"github.com/uphy/commandbeat/parser"
)

type commandSpec interface {
	command.Spec
	parser() parser.Parser
	name() string
	debug() bool
}

// defaultCommandSpec implements command.Spec interface.
type defaultCommandSpec struct {
	_name   string
	command string
	args    []string
	_parser parser.Parser
	_debug  bool
}

func (cb *Commandbeat) newSpec(name string, task *config.TaskConfig) (command.Spec, error) {
	parser, err := task.Parser()
	if err != nil {
		return nil, err
	}

	if task.Shell {
		script, err := task.Command()
		if err != nil {
			return nil, err
		}
		s, err := cb.scriptManager.createScript(name, script[0])
		if err != nil {
			return nil, fmt.Errorf("Failed to create script file. (dir=%s, name=%s, err=%s)", cb.scriptManager.directory, name, err)
		}
		return &defaultCommandSpec{
			name,
			s.Command(),
			s.Args(),
			parser,
			task.Debug,
		}, nil
	}
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
	return &defaultCommandSpec{name, commandName, commandArgs, parser, task.Debug}, nil
}

func (c *defaultCommandSpec) Command() string {
	return c.command
}

func (c *defaultCommandSpec) Args() []string {
	return c.args
}

func (c *defaultCommandSpec) parser() parser.Parser {
	return c._parser
}

func (c *defaultCommandSpec) name() string {
	return c._name
}

func (c *defaultCommandSpec) debug() bool {
	return c._debug
}
