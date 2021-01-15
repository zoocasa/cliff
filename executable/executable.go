package executable

import "github.com/zoocasa/cliff/cliff"

// Execute configures the CLI and executes the root command
func Execute() {
	cliff.ConfigureFromFile("cli.yml")
	cliff.Execute()
}
