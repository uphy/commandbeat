package command

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"
)

type testCommand struct {
	command string
	args    []string
}

func (t *testCommand) Command() string {
	return t.command
}
func (t *testCommand) Args() []string {
	return t.args
}
func (t *testCommand) create() *exec.Cmd {
	return exec.Command(t.command, t.args...)
}

type testHandler struct {
	beforeStartErr error
	afterExitErr   error
	status         int
	stdoutErr      error
	stderrErr      error
	handledErr     error
	stdout         string
	stderr         string
}

func (t *testHandler) BeforeStart(spec Spec) error {
	return t.beforeStartErr
}
func (t *testHandler) HandleStdOut(spec Spec, out string) error {
	t.stdout = out
	return t.stdoutErr
}
func (t *testHandler) HandleStdErr(spec Spec, err string) error {
	t.stderr = err
	return t.stderrErr
}
func (t *testHandler) HandleErr(spec Spec, err error) {
	t.handledErr = err
}
func (t *testHandler) AfterExit(spec Spec, status int) error {
	t.status = status
	return t.afterExitErr
}

func TestRunnerRun(t *testing.T) {
	handler := &testHandler{}
	runner := NewRunner(handler)
	spec := &testCommand{"echo", []string{"hello"}}
	runner.Run(spec)
	if handler.handledErr != nil {
		t.Error(handler.handledErr)
	}
}

func TestRunnerRunBeforeError(t *testing.T) {
	handler := &testHandler{beforeStartErr: errors.New("err")}
	runner := NewRunner(handler)
	spec := &testCommand{"echo", []string{"hello"}}
	runner.Run(spec)
	if handler.handledErr == nil {
		t.Error("error not retuned.")
	}
}

func TestRunnerRunStartError(t *testing.T) {
	handler := &testHandler{}
	runner := NewRunner(handler)
	spec := &testCommand{"echoaaaaaaaaaaaaaaaaaaa", []string{"hello"}}
	runner.Run(spec)
	if handler.handledErr == nil {
		t.Error("error not retuned.")
	}
}

func TestRunnerRunAfterExitError(t *testing.T) {
	handler := &testHandler{afterExitErr: errors.New("err")}
	runner := NewRunner(handler)
	spec := &testCommand{"echo", []string{"hello"}}
	runner.Run(spec)
	if handler.handledErr == nil {
		t.Error("error not retuned.")
	}
}

func TestStartCommand(t *testing.T) {
	spec := &testCommand{"echo", []string{"hello"}}
	cmd := spec.create()
	if err := startCommand(spec, cmd, &testHandler{}); err != nil {
		t.Error(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Error(err)
	}
}

func TestStartCommandStdErr(t *testing.T) {
	spec := &testCommand{"sh", []string{"-c", "echo hello 1>&2"}}
	cmd := spec.create()
	if err := startCommand(spec, cmd, &testHandler{}); err != nil {
		t.Error(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Error(err)
	}
}

func TestStartCommandErr(t *testing.T) {
	spec := &testCommand{"false", []string{}}
	cmd := spec.create()
	if err := startCommand(spec, cmd, &testHandler{}); err != nil {
		t.Error(err)
	}
	if err := cmd.Wait(); err == nil {
		t.Error("expected error")
	}
}

func TestStartCommandStdOutErr(t *testing.T) {
	spec := &testCommand{"echo", []string{"hello"}}
	cmd := spec.create()
	h := &testHandler{
		stdoutErr: errors.New("err"),
	}
	if err := startCommand(spec, cmd, h); err != nil {
		t.Error(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Error(err)
	}
	if h.handledErr == nil {
		t.Error("error not handled.")
	}
}

func TestStartCommandNoSuchCommand(t *testing.T) {
	spec := &testCommand{"echoaaaaaa", []string{"hello"}}
	cmd := spec.create()
	if err := startCommand(spec, cmd, &testHandler{}); err == nil {
		t.Error("expected error")
	}
}

type closingBuffer struct {
	*bytes.Buffer
}

func (c *closingBuffer) Close() (err error) {
	return nil
}
func newBufferString(s string) io.ReadCloser {
	return &closingBuffer{bytes.NewBufferString(s)}
}

func TestHandleStream(t *testing.T) {
	var got string
	handleStream(newBufferString("oneline"), func(line string) {
		got = line
	}, func(err error) {
		t.Error(err)
	})
	if got != "oneline" {
		t.Errorf("unexpected line: %s", got)
	}
}

type errReader struct {
	io.ReadCloser
}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("err")
}

func TestHandleStreamError(t *testing.T) {
	var got error
	handleStream(&errReader{newBufferString("oneline")}, func(line string) {

	}, func(err error) {
		got = err
	})
	if got == nil {
		t.Error("error not returned")
	}
}

func TestExtractExitStatus(t *testing.T) {
	status, err := extractExitStatus(&exec.ExitError{
		ProcessState: &os.ProcessState{},
	})
	if err != nil {
		t.Error(err)
	}
	if status != 0 {
		t.Error("unexpected status: ", status)
	}

	status, err = extractExitStatus(errors.New(""))
	if err == nil {
		t.Error("error nil")
	}
}
