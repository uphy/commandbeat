package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type (
	// Runner runs command.  command result is handled with Handler interface.
	Runner struct {
		handler Handler
	}
	// Spec represents command and its arguments.
	Spec interface {
		Command() string
		Args() []string
	}
	// Handler handles command result.
	Handler interface {
		BeforeStart(spec Spec) error
		HandleStdOut(spec Spec, out string) error
		HandleStdErr(spec Spec, err string) error
		HandleErr(spec Spec, err error)
		AfterExit(spec Spec, status int) error
	}
)

// NewRunner creates a new runner.
func NewRunner(handler Handler) *Runner {
	return &Runner{handler}
}

// Run runs the specified command.
func (r *Runner) Run(spec Spec) {
	if err := r.handler.BeforeStart(spec); err != nil {
		r.handler.HandleErr(spec, err)
		return
	}
	cmd := exec.Command(spec.Command(), spec.Args()...)
	if err := startCommand(spec, cmd, r.handler); err != nil {
		r.handler.HandleErr(spec, err)
		return
	}
	statusCode := 0
	if err := cmd.Wait(); err != nil {
		status, err2 := extractExitStatus(err)
		if err2 != nil {
			r.handler.HandleErr(spec, fmt.Errorf("failed to wait for command exit. (cmd=%v, err=%v)", cmd.Args, err))
			return
		}
		statusCode = status
	}
	if err := r.handler.AfterExit(spec, statusCode); err != nil {
		r.handler.HandleErr(spec, fmt.Errorf("failed to process command exit. (cmd=%v, err=%v)", cmd.Args, err))
		return
	}
}

func startCommand(spec Spec, cmd *exec.Cmd, handler Handler) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout. (cmd=%v, err=%s)", cmd.Args, err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdout.Close()
		return fmt.Errorf("failed to get stderr. (cmd=%v, err=%s)", cmd.Args, err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command. (cmd=%v, err=%s)", cmd.Args, err)
	}

	errHandler := func(err error) {
		if err != nil {
			handler.HandleErr(spec, err)
		}
	}
	go handleStream(stdout, func(line string) {
		errHandler(handler.HandleStdOut(spec, line))
	}, errHandler)
	go handleStream(stderr, func(line string) {
		errHandler(handler.HandleStdErr(spec, line))
	}, errHandler)
	return nil
}

func handleStream(reader io.ReadCloser, consumer func(line string), errHandler func(err error)) {
	defer reader.Close()
	bufReader := bufio.NewReader(reader)
	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			// FIXME
			if perr, ok := err.(*os.PathError); ok && strings.Contains(perr.Error(), "file already closed") {
				return
			}
			errHandler(fmt.Errorf("failed to read stdout or stderr stream: %v", err))
			return
		}
		consumer(string(line))
	}
}

func extractExitStatus(err error) (int, error) {
	if e2, ok := err.(*exec.ExitError); ok {
		if s, ok := e2.Sys().(syscall.WaitStatus); ok {
			exitStatus := s.ExitStatus()
			return exitStatus, nil
		}
		return -1, errors.New("can not get exit status code in your environment")
	}
	return -1, err
}
