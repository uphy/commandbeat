package beater

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"syscall"
	"time"

	"github.com/uphy/commandbeat/parser"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"
)

type (
	commandRunner struct {
		client beat.Client
	}
	commandSpec struct {
		name       string
		command    string
		args       []string
		exitStatus int
		debug      bool
	}
)

func newCommand(name string, commandName string, debug bool, args ...string) *commandSpec {
	cmd := commandSpec{name, commandName, args, 0, debug}
	return &cmd
}

func newCommandRunner(client beat.Client) *commandRunner {
	return &commandRunner{client}
}

func (c *commandRunner) debug(spec *commandSpec, msg string, args ...interface{}) {
	if spec.debug {
		logp.Info("[%s] %s", spec.name, fmt.Sprintf(msg, args...))
	}
}

func (c *commandRunner) run(spec *commandSpec, parser parser.Parser) {
	c.debug(spec, "executing command(command=%s, args=%v)", spec.command, spec.args)
	cmd := exec.Command(spec.command, spec.args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logp.Err("failed to get stdout. (cmd=%v, err=%s)", cmd.Args, err)
	}
	stdoutReader := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		logp.Err("failed to start command. (cmd=%v, err=%s)", cmd.Args, err)
		return
	}
	defer stdout.Close()
	for {
		line, _, err := stdoutReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			logp.Err("failed to read stdin. (cmd=%v, err=%s)", cmd.Args, err)
		}
		lineStr := string(line)
		c.debug(spec, "<stdout>%s", line)
		v, err := parser.Parse(lineStr)
		if err != nil {
			logp.Err("failed to parse the line got from stdin. (cmd=%v, line=%s, err=%s)", cmd.Args, lineStr, err)
			continue
		}
		c.debug(spec, "<parsed>%#v", v)
		var timestamp time.Time
		if t, ok := v["@timestamp"]; ok {
			timestamp = t.(time.Time)
			delete(v, "@timestamp")
		} else {
			timestamp = time.Now()
		}
		v["type"] = spec.name
		event := beat.Event{
			Timestamp: timestamp,
			Fields:    v,
		}
		c.debug(spec, "<event>%#v", event)
		if !spec.debug {
			c.client.Publish(event)
		}
	}
	if err := cmd.Wait(); err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				exitStatus := s.ExitStatus()
				if exitStatus != spec.exitStatus {
					logp.Err("Unexpected exit status. (cmd=%v, status=%d)", cmd.Args, exitStatus)
				}
			} else {
				panic(errors.New("can not get exit status code in your environment"))
			}
		} else {
			logp.Err("failed to wait for command exit. (cmd=%v, err=%s)", cmd.Args, err)
		}
	}
}
