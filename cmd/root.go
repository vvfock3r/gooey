package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vvfock3r/gooey/kernel/load"
	"github.com/vvfock3r/gooey/kernel/module/mysql"
)

var rootCmd = &cobra.Command{
	Use:           "gooey",
	Short:         "Simple Command-Line Interface Template\nFor details, please refer to https://github.com/vvfock3r/gooey",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		for _, m := range load.ModuleList {
			m.MustCheck(cmd)
		}

		for _, m := range load.ModuleList {
			err = m.Initialize(cmd)
			if err != nil {
				return err
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		//for {
		//	now := time.Now().Format(time.DateTime)
		//	logger.Debug(now)
		//	logger.Info(now)
		//	logger.Warn(now)
		//	logger.Error(now)
		//	fmt.Println()
		//	time.Sleep(time.Second)
		//}
		var v string
		mysql.DB.Get(&v, "select @@version")
		fmt.Println(v)
	},
}

func init() {
	// register flags or others
	for _, m := range load.ModuleList {
		m.Register(rootCmd)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
