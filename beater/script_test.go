package beater

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func testInTmpDir(t *testing.T, f func(tmpDir string)) {
	tmp, err := ioutil.TempDir("", "commandbeattest")
	if err != nil {
		t.Error(err)
		return
	}
	f(tmp)
	if err := os.RemoveAll(tmp); err != nil {
		t.Error(err)
	}
}

func TestScriptManagerNewAndDelete(t *testing.T) {
	testInTmpDir(t, func(tmp string) {
		dir := filepath.Join(tmp, "scripts")

		scriptManager, err := newScriptManager(dir)
		if scriptManager == nil {
			t.Error(err)
		}
	})
}

func TestScriptManagerNewErrNotCreated(t *testing.T) {
	testInTmpDir(t, func(tmp string) {
		dir := filepath.Join(tmp, "scripts")
		// create file
		f, err := os.Create(dir)
		if err != nil {
			t.Error(err)
			return
		}
		defer f.Close()

		// already exists as file.  we should not be create the directory.
		if _, err := newScriptManager(dir); err == nil {
			t.Error("error not returned. special character file should not be created.")
		}
	})
}

func TestScriptManagerCreate(t *testing.T) {
	testInTmpDir(t, func(tmp string) {
		scriptManager, _ := newScriptManager(tmp)
		script, err := scriptManager.createScript("script", `#!/bin/bash
echo hello
`)
		if err != nil {
			t.Error(err)
			return
		}
		if f, _ := os.Stat(script.path); f.Mode() != 0700 {
			t.Error(err)
			return
		}
	})
}

func TestScriptManagerOverwrite(t *testing.T) {
	testInTmpDir(t, func(tmp string) {
		scriptManager, _ := newScriptManager(tmp)
		script, _ := scriptManager.createScript("script", `#!/bin/bash
echo hello
`)
		script, err := scriptManager.createScript("script", `#!/bin/bash
	echo hello
`)
		if err != nil {
			t.Error(err)
			return
		}
		if f, _ := os.Stat(script.path); f.Mode() != 0700 {
			t.Error(err)
			return
		}
	})
}

func TestScriptManagerClean(t *testing.T) {
	testInTmpDir(t, func(tmp string) {
		scriptManager, _ := newScriptManager(tmp)
		scriptManager.createScript("script", `#!/bin/bash
	echo hello
`)
		if err := scriptManager.clean(); err != nil {
			t.Error(err)
		}
		if _, err := os.Stat(tmp); !os.IsNotExist(err) {
			t.Error("not cleaned up.")
		}
	})
}

func TestScriptManagerScriptCommandAndArgs(t *testing.T) {
	testInTmpDir(t, func(tmp string) {
		scriptManager, _ := newScriptManager(tmp)
		script, _ := scriptManager.createScript("script", `#!/bin/bash
	echo hello
`)
		if script.Command() != script.path {
			t.Error("invalid command")
		}
		if len(script.Args()) != 0 {
			t.Error("invalid args")
		}
	})
}
