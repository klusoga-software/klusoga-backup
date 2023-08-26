package main

import (
	"github.com/klusoga-software/klusoga-backup-agent/cmd"
	"github.com/klusoga-software/klusoga-backup-agent/pkg/build"
)

var version = "Dev"

func main() {
	build.Version = version
	cmd.Execute()
}
