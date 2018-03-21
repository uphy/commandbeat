package beater

import (
	"errors"
	"fmt"

	"github.com/uphy/commandbeat/command"
	"github.com/uphy/commandbeat/config"
	"github.com/uphy/commandbeat/parser"
)

// defaultCommandSpec implements command.Spec interface.
type commandSpec struct {
	name    string
	command string
	args    []string
	parser  parser.Parser
	debug   bool
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
		return &commandSpec{
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
	return &commandSpec{name, commandName, commandArgs, parser, task.Debug}, nil
}

func (c *commandSpec) Command() string {
	return c.command
}

func (c *commandSpec) Args() []string {
	return c.args
}
