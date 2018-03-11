// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"fmt"

	"github.com/mattn/go-shellwords"
	"github.com/uphy/commandbeat/parser"
)

type (
	// Config is the root struct for the config file.
	Config struct {
		Tasks map[string]TaskConfig `config:"tasks"`
	}
	// TaskConfig represents the task config.
	TaskConfig struct {
		CommandRaw       interface{}            `config:"command"`
		WorkingDirectory string                 `config:"workdir"`
		ParserRaw        map[string]interface{} `config:"parser"`
		Schedule         string                 `config:"schedule"`
	}
)

var (
	// DefaultConfig is the default config.
	DefaultConfig = Config{
		Tasks: map[string]TaskConfig{},
	}
)

// Command parses the command line as array.
func (t *TaskConfig) Command() ([]string, error) {
	if array, ok := t.CommandRaw.([]interface{}); ok {
		var v = []string{}
		for _, elm := range array {
			if s, ok := elm.(string); ok {
				v = append(v, s)
			} else {
				return nil, fmt.Errorf("command argument must be a string. (argument=%v)", elm)
			}
		}
		return v, nil
	}
	if s, ok := t.CommandRaw.(string); ok {
		command, err := shellwords.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse command line. (commandline=%s)", s)
		}
		return command, nil
	}
	return nil, fmt.Errorf("unsupported type of command line. (command=%v)", t.CommandRaw)
}

// Parser creates parser from config.
func (t *TaskConfig) Parser() (parser.Parser, error) {
	if t.ParserRaw == nil {
		return parser.NewMessageParser(), nil
	}
	return parser.NewParser(t.ParserRaw)
}
