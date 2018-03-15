package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"syscall"

	"github.com/uphy/commandbeat/parser"

	"github.com/elastic/beats/libbeat/logp"
)

type (
	CommandRunner struct {
		handler Handler
	}
	Spec struct {
		Name       string
		Command    string
		Args       []string
		ExitStatus int
		Parser     parser.Parser
		Debug      bool
	}
	Handler interface {
		HandleStdOut(spec *Spec, out string) error
		HandleStdErr(spec *Spec, err string) error
	}
)

func NewCommand(name string, commandName string, parser parser.Parser, debug bool, args ...string) *Spec {
	cmd := Spec{name, commandName, args, 0, parser, debug}
	return &cmd
}

func NewCommandRunner(handler Handler) *CommandRunner {
	return &CommandRunner{handler}
}

func (c *Spec) LogDebug(msg string, args ...interface{}) {
	if c.Debug {
		logp.Info("[%s] %s", c.Name, fmt.Sprintf(msg, args...))
	}
}

func (c *CommandRunner) Run(spec *Spec) error {
	cmd := exec.Command(spec.Command, spec.Args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout. (cmd=%v, err=%s)", cmd.Args, err)
	}
	stdoutReader := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command. (cmd=%v, err=%s)", cmd.Args, err)
	}
	defer stdout.Close()
	for {
		line, _, err := stdoutReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read stdin. (cmd=%v, err=%s)", cmd.Args, err)
		}
		if err := c.handler.HandleStdOut(spec, string(line)); err != nil {
			return err
		}
	}
	if err := cmd.Wait(); err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				exitStatus := s.ExitStatus()
				if exitStatus != spec.ExitStatus {
					return fmt.Errorf("Unexpected exit status. (cmd=%v, status=%d)", cmd.Args, exitStatus)
				}
			} else {
				return errors.New("can not get exit status code in your environment")
			}
		} else {
			return fmt.Errorf("failed to wait for command exit. (cmd=%v, err=%s)", cmd.Args, err)
		}
	}
	return nil
}
