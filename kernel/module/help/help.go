package help

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Help implement the Module interface
type Help struct {
	HiddenShortFlag   bool
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

	if h.HiddenShortFlag {
		cmd.PersistentFlags().Bool("help", false, "help message")
	} else {
		cmd.PersistentFlags().BoolP("help", "h", false, "help message")
	}

	if h.HiddenHelpCommand {
		cmd.SetHelpCommand(&cobra.Command{Use: "no-help", Hidden: true})
	}
}

func (h *Help) MustCheck(*cobra.Command) {}

func (h *Help) Initialize(*cobra.Command) error {
	return nil
}
