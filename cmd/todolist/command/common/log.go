package common

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
)

// ParseCmdLog retrieves log parameters
func ParseCmdLog(cmd *cobra.Command) (int, string, error) {
	verbosity, err := strconv.Atoi(cmd.Flag("v").Value.String())
	if err != nil {
		return 0, "", errors.New("verbosity must be an integer value")
	}
	logFormatter := cmd.Flag("log-formatter").Value.String()

	return verbosity, logFormatter, nil
}
