package note

import (
	"errors"
	"github.com/spf13/cobra"
)

var (
	pageSize   int
	serverName string
)

var NoteCmd = &cobra.Command{
	Use:     "note",
	Aliases: []string{"notes"},
	Short:   "Manage and view server notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("subcommand error")
	},
}

func init() {
	NoteCmd.AddCommand(noteListCmd)
	NoteCmd.AddCommand(noteCreateCmd)
	NoteCmd.AddCommand(noteDeleteCmd)
}
