//version
package main

import (
	"fmt"
	"runtime"
)

var myVersion = "0.0.1"

var cmdVersion = &Command{
	Run:               runVersion,
	UsageLine:         "version",
	Short:             "print screen-dump-service version",
	Long:              `Version prints the screen-dump-service version`,
	DefaultParameters: "",
}

func runVersion(cmd *Command, args []string) bool {
	if len(args) != 0 {
		cmd.Usage()
	}

	fmt.Printf("version %s %s %s\n", myVersion, runtime.GOOS, runtime.GOARCH)
	return true
}
