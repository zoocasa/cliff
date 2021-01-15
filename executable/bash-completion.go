package executable

import (
	"github.com/zoocasa/cliff/cliff"
)

func init() {
	cliff.AddRunToCommand("bash-completion", func(c *cliff.Command) {
		outputPath := c.FlagString("output")
		cliff.ConfigureAndGenerateBashCompletionFile(outputPath)
	})
}
