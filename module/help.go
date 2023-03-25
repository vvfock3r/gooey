package module

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Help implement the Module interface
type Help struct {
	HiddenHelpCommand bool
}

func (h *Help) Register(cmd *cobra.Command) {
	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		if command.Long != "" {
			fmt.Printf("%s\n\n", command.Long)
		} else {
			fmt.Printf("%s\n\n", command.Short)
		}
		fmt.Printf("%s", command.UsageString())
		os.Exit(0)
	})

	if h.HiddenHelpCommand {
		cmd.SetHelpCommand(&cobra.Command{Use: "no-help", Hidden: true})
	}

	cmd.PersistentFlags().BoolP("help", "h", false, "help message")
}

func (h *Help) MustCheck(*cobra.Command) {}

func (h *Help) Initialize(*cobra.Command) error {
	return nil
}
