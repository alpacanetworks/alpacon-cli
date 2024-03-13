package note

import (
	"github.com/alpacanetworks/alpacon-cli/api/note"
	"github.com/alpacanetworks/alpacon-cli/client"
	"github.com/alpacanetworks/alpacon-cli/utils"
	"github.com/spf13/cobra"
)

var noteDeleteCmd = &cobra.Command{
	Use:   "delete [NOTE ID]",
	Short: "Delete a specified note",
	Long: `
	This command permanently deletes a specified note from the Alpacon server. 
	It's important to verify that you have the necessary permissions to delete a note before using this command. 
	The command requires an exact note ID as its argument.
	`,
	Example: ` 
	alpacon server delete [NOTE ID]	
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		noteID := args[0]

		alpaconClient, err := client.NewAlpaconAPIClient()
		if err != nil {
			utils.CliError("Connection to Alpacon API failed: %s. Consider re-logging.", err)
		}

		err = note.DeleteNote(alpaconClient, noteID)
		if err != nil {
			utils.CliError("Failed to delete the note with ID %s. Error: %s. Please check the note ID and your permissions, and try again.", noteID, err)
		}

		utils.CliInfo("Note successfully deleted: %s.", noteID)
	},
}
