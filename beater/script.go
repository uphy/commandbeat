package beater

import (
	"os"
	"path/filepath"

	"io/ioutil"

	"github.com/hashicorp/go-multierror"
)

type scriptManager struct {
	directory   string
	scriptFiles []scriptFile
}

func newScriptManager(directory string) (*scriptManager, error) {
	if err := os.MkdirAll(directory, 0700); err != nil {
		return nil, err
	}
	return &scriptManager{directory, []scriptFile{}}, nil
}

func (s *scriptManager) createScript(name string, script string) (*scriptFile, error) {
	file := filepath.Join(s.directory, name)
	if err := ioutil.WriteFile(file, []byte(script), 0700); err != nil {
		return nil, err
	}
	scriptFile := scriptFile{file}
	s.scriptFiles = append(s.scriptFiles, scriptFile)
	return &scriptFile, nil
}

func (s *scriptManager) clean() error {
	var err error
	for _, script := range s.scriptFiles {
		if err := script.delete(); err != nil {
			err = multierror.Append(err, err)
		}
	}
	if err := os.Remove(s.directory); err != nil {
		err = multierror.Append(err)
	}
	return err
}

type scriptFile struct {
	path string
}

func (s *scriptFile) Command() string {
	return s.path
}

func (s *scriptFile) Args() []string {
	return []string{}
}

func (s *scriptFile) delete() error {
	return os.Remove(s.path)
}
