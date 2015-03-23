// ScreenDump
package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var IsDebug *bool
var exitStatus = 0

var commands = []*Command{
	cmdVersion,
	cmdServer,
}

func main() {
	glog.MaxSize = 1024 * 1024 * 32
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		for _, cmd := range commands {
			if len(args) >= 2 && cmd.Name() == args[1] && cmd.Run != nil {
				fmt.Fprintf(os.Stderr, "Default Parameters:%s\n", cmd.DefaultParameters)
				cmd.Flag.PrintDefaults()
			}
		}
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			IsDebug = cmd.IsDebug
			if !cmd.Run(cmd, args) {
				fmt.Fprintf(os.Stderr, "\n")
				cmd.Flag.Usage()
				fmt.Fprintf(os.Stderr, "Default Parameters:%d\n", cmd.DefaultParameters)
				cmd.Flag.PrintDefaults()
			}
			exit()
			return
		}
	}

	fmt.Fprintf(os.Stderr, "screen-dump-service: unknown subcommand %q\nRun 'screen-dump-service help' for usage.\n", args[0])
	exit()
}

var usageTemplate = `
Usage:

	screen-dump-service command [arguments]

The commands are:
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "screen-dump-service help [command]" for more information about a command.

`

var helpTemplate = `{{if .Runnable}}Usage: screen-dump-service {{.UsageLine}}{{end}}
{{.Long}}
`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, commands)
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

// help implements the 'help' command.
func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'weed help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: screen-dump-service help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'weed help'
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'weed help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'screen-dump-service help'.\n", arg)
	os.Exit(2) // failed at 'weed help cmd'
}

var atexitFuncs []func()

func exit() {
	for _, f := range atexitFuncs {
		f()
	}
	os.Exit(exitStatus)
}
