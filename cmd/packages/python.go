package packages

import (
	"errors"
	"github.com/spf13/cobra"
)

var pythonCmd = &cobra.Command{
	Use:   "python",
	Short: "Python packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("subcommand error")
	},
}

func init() {
	pythonCmd.AddCommand(pythonPackageListCmd)
	pythonCmd.AddCommand(pythonPackageUploadCmd)
	pythonCmd.AddCommand(pythonPackageDownloadCmd)
}
