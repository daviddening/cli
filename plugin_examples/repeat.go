/**
* This is an example plugin where we use both arguments and flags. The plugin
* will echo all arguments passed to it. The flag -uppercase will upcase the
* arguments passed to the command. The help flag will print the usage text for
* this command and exit, ignoring any other arguments passed.
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

type RepeatParams struct {
	help   *bool
	dryrun *bool
}

var permitted_cmds = []string{
	"app",
	"bind-service",
	"delete",
	"e",
	"env",
	"events",
	"files",
	"logs",
	"map-route",
	"push",
	"restage",
	"restart",
	"scale",
	"se",
	"set-env",
	"stacks",
	"start",
	"stop",
	"unbind-service",
	"unmap-route",
	"unset-env"}

func main() {
	plugin.Start(new(RepeatParams))
}

func (repeatParams *RepeatParams) Run(cliConnection plugin.CliConnection, args []string) {
	// Initialize flags
	echoFlagSet := flag.NewFlagSet("echo", flag.ExitOnError)
	help := echoFlagSet.Bool("help", false, "passed to display help text")
	dryrun := echoFlagSet.Bool("dryrun", false, "passed to see what commands would run")

	// Parse starting from [1] because the [0]th element is the
	// name of the command
	err := echoFlagSet.Parse(args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *help {
		printHelp()
		os.Exit(0)
	}

	var itemToEcho string
	cmd := echoFlagSet.Args()[0]
	var apps_string string = echoFlagSet.Args()[1]
	cmd_args := echoFlagSet.Args()[2:]

	if !contains(permitted_cmds, cmd) {
		fmt.Println("Please use a permitted command: \n" + strings.Join(permitted_cmds, "\n"))
		os.Exit(1)
	}

	apps_string = strings.TrimPrefix(apps_string, "[")
	apps_string = strings.TrimSuffix(apps_string, "]")

	apps := strings.Split(apps_string, ",")

	for i := 0; i < len(apps); i++ {
		itemToEcho = cmd + " " + apps[i] + " " + strings.Join(cmd_args, " ")

		if *dryrun {
			fmt.Println("Repeat will run \"cf " + itemToEcho + "\"")
			continue
		}

		fmt.Println("Starting to run \"" + itemToEcho + "\" ...")

		commands := make([]string, len(cmd_args)+2)
		commands[0] = cmd
		commands[1] = apps[i]
		copy(commands[2:], cmd_args[0:])

		output, err := cliConnection.CliCommand(commands...)
		if err != nil {
			fmt.Println("PLUGIN ERROR: Error from CliCommand: ", err)
		} else {
			fmt.Println("OUTPUT: \n", output)
		}
		fmt.Println("\n")
	}
}

func (repeatParams *RepeatParams) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Repeat",
		Commands: []plugin.Command{
			{
				Name:     "repeat",
				HelpText: "Repeat a command for an array of apps. To obtain more information use --help",
			},
		},
	}
}

func printHelp() {
	fmt.Println(`
	cf repeat CMD APPS ARGS

	OPTIONAL PARAMS:
	-help: used to display this additional output.
	-dryrun: see what commands would run without running them.

	REQUIRED PARAMS:
	cmd : command to repeat (eg. bind-service) 
	apps: apps to run command on without spaces  (eg. [production-blue,production-green] )
	args: additional arguments in quotes (eg. 'my-elephant-sql' )

	EXAMPLE:
	cf repeat map-route [production-blue,production-green] cfapps.io -n donottarget
	`)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
