package main

import "mcolomerc/vvp2cli/cmd"

var (
	version   = "dev"
	commit    = "none"
	buildTime = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, buildTime)
	cmd.Execute()
}
