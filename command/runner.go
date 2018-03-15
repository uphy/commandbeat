package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

type (
	Runner struct {
		handler Handler
	}
	Spec interface {
		Command() string
		Args() []string
	}
	Handler interface {
		BeforeStart(spec Spec) error
		HandleStdOut(spec Spec, out string) error
		HandleStdErr(spec Spec, err string) error
		AfterExit(spec Spec, status int) error
	}
)

func NewRunner(handler Handler) *Runner {
	return &Runner{handler}
}

func (c *Runner) Run(spec Spec) error {
	if err := c.handler.BeforeStart(spec); err != nil {
		return err
	}
	cmd := exec.Command(spec.Command(), spec.Args()...)
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
				if err := c.handler.AfterExit(spec, exitStatus); err != nil {
					return err
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
