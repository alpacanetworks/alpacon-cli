package authority

import (
	"errors"
	"github.com/spf13/cobra"
)

var AuthorityCmd = &cobra.Command{
	Use:   "authority",
	Short: "Commands to manage and interact with certificate authorities",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("subcommand error")
	},
}

func init() {
	AuthorityCmd.AddCommand(authorityCreateCmd)
	AuthorityCmd.AddCommand(authorityListCmd)
	AuthorityCmd.AddCommand(authorityDetailCmd)
	AuthorityCmd.AddCommand(authorityDownloadCmd)
}
