package command

import (
	"fmt"

	ver "github.com/odacremolbap/rest-demo/pkg/version"
	"github.com/spf13/cobra"
)

// VersionCmd shows version information
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:\t%s\nDate (UTC):\t%s\n", ver.Version, ver.Date)
	},
}
