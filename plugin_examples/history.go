/**
* This is an example plugin where we use both arguments and flags. The plugin
* will echo all arguments passed to it. The flag -uppercase will upcase the
* arguments passed to the command. The help flag will print the usage text for
* this command and exit, ignoring any other arguments passed.
 */
package main

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
)

type BasicPlugin struct{}

func main() {
	plugin.Start(new(BasicPlugin))
}

func (c *BasicPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	var itemToEcho string

	itemToEcho = ` 
	Summer 2010  - Started by VMWare
	April  2011  - cloudfoundry.com goes live!
	April  2012  - BOSH released
	Fall   2012  - Pivotal Labs engagement
	April  2013  - Pivotal acquires CF 
	`


	fmt.Println(itemToEcho)
}

func (c *BasicPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "History",
		Commands: []plugin.Command{
			{
				Name:     "history",
				HelpText: "Display a brief history of the CF project.",
			},
		},
	}
}

func printHelp() {
	fmt.Println(`
cf history
		`)
}
