package cliff

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type run struct {
	Name string
	Run  func(c Command, args []string)
}

var rootCmd = &Command{}
var commands = map[string]*Command{}
var runs = map[string]func(c *Command){}

func init() {
	log.SetFlags(0)
}

// Configure sets the content of the yaml config file and sets up the commands
func Configure(yamlConfigContent []byte) {
	setupRootCmd(yamlConfigContent)
	attachRunToCommands()
}

// ConfigureFromFile reads the contented of a passed file path and then calls Configure with it
func ConfigureFromFile(path string) error {
	yamlConfigContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	Configure(yamlConfigContent)
	return nil
}

// ConfigureSubcommandFromFile reads a file, generates a command and attaches it to the root command as a subcommand
func ConfigureSubcommandFromFile(path string) error {
	yamlConfigContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	cmd := commandFromConfigFile(yamlConfigContent)
	if rootCmd.Name != cmd.Name {
		(*rootCmd.cobraCmd).AddCommand(cmd.buildCommand().cobraCmd)
	}
	return nil
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := rootCmd.cobraCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// ConfigureAndExecute runs ConfigureFromFile with default
// locations for the yaml config file and then runs Execute
func ConfigureAndExecute() {
	ConfigureFromFile("cli.yml")
	Execute()
}

// AddRunToCommand provies a mechanism to attach a Run function to a command
func AddRunToCommand(name string, runFunc func(c *Command)) {
	runs[name] = runFunc
}

func setupRootCmd(config []byte) {
	rootCmd = commandFromConfigFile(config).buildCommand()
	addVerboseFlagToRootCmd()
	rootCmd.cobraCmd.SetHelpCommand(&cobra.Command{}) // Remove default help subcommand
}

func commandFromConfigFile(config []byte) *Command {
	var command Command
	if err := yaml.Unmarshal(config, &command); err != nil {
		log.Fatal(err)
	}
	return &command
}

func attachRunToCommands() {
	for name := range runs {
		if c, ok := commands[name]; ok {
			c.cobraCmd.Run = func(_ *cobra.Command, args []string) {
				c.args = args // set args before runing the function
				runs[c.Name](c)
				c.args = nil // Clear the args as they should not be part of the permanent command data
			}
		}
	}
}
