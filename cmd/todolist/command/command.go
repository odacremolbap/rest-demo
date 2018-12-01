package command

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/odacremolbap/rest-demo/cmd/todolist/command/common"
	"github.com/odacremolbap/rest-demo/cmd/todolist/command/server"
	"github.com/spf13/cobra"
)

var (
	logFormatter string

	// TODOCmd is the base command for the binary
	TODOCmd = &cobra.Command{
		Use:   "todolist",
		Short: "Todo",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			verbosity, logFormatter, err := common.ParseCmdLog(cmd)
			if err != nil {
				logrus.Error(err)
				_ = cmd.Usage()
				os.Exit(-1)
			}
			// configure logs for all subcommands
			common.ConfigureLogrusLogger(logFormatter, verbosity)
		},
		Run: runHelp,
	}
)

func init() {

	TODOCmd.PersistentFlags().IntP("v", "v", 1, "verbosity level")
	TODOCmd.PersistentFlags().StringP("log-formatter", "", "server", "one of server|text|json")

}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

// Execute subcommand to load logs
func Execute() {
	TODOCmd.SilenceUsage = true
	TODOCmd.AddCommand(server.ServerCmd)
	TODOCmd.AddCommand(VersionCmd)

	if err := TODOCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
