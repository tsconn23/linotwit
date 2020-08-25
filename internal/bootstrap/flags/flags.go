package flags

import (
	"flag"
	"fmt"
	"os"
)

const (
	defaultConfigDirectory = "./cfg"
	defaultConfigFile = "configuration.toml"
)

// Common is an interface that defines AP for the common command-line flags
type Common interface {
	ConfigDirectory() string
	ConfigFileName() string
	Parse([]string)
	Help()
}

//Default contains a common set of flags
type Default struct {
	FlagSet           *flag.FlagSet
	additionalUsage   string
	configDirectory   string
	configFileName    string
}

// New returns a Default struct with an empty additional usage string.
func New() *Default {
	return &Default {
		FlagSet: flag.NewFlagSet("", flag.ExitOnError),
	}
}

// Parse parses the passed in command-lie arguments looking to the default set of common flags
func (d *Default) Parse(arguments []string) {
	d.FlagSet.StringVar(&d.configFileName, "cfgFile", defaultConfigFile, "")
	d.FlagSet.StringVar(&d.configDirectory, "cfgDir", defaultConfigDirectory, "")

	d.FlagSet.Usage = d.helpCallback

	err := d.FlagSet.Parse(arguments)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// ConfigDirectory returns the directory where the config file(s) are located, if it was specified.
func (d *Default) ConfigDirectory() string {
	return d.configDirectory
}

// ConfigFileName returns the name of the local configuration file
func (d *Default) ConfigFileName() string {
	return d.configFileName
}

// Help displays the usage help message and exit.
func (d *Default) Help() {
	d.helpCallback()
}

// commonHelpCallback displays the help usage message and exits
func (d *Default) helpCallback() {
	fmt.Printf(
		"Usage: %s [options]\n"+
			"Server Options:\n"+
			"    --cfgFile <name>                Indicates name of the local configuration file. Defaults to configuration.toml\n"+
			"    --cfgDir                        Specify local configuration directory\n"+
			"%s\n"+
			"Common Options:\n"+
			"	-h, --help                      Show this message\n",
		os.Args[0], d.additionalUsage,
	)
	os.Exit(0)
}