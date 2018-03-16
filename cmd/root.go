package cmd

import (
	"github.com/uphy/commandbeat/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"
)

// Name of this beat
var Name = "commandbeat"

// Version of this beat
var Version = "0.0.2pre"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmd(Name, Version, beater.New)
