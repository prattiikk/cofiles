package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var spaceLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all spaces",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all spaces...")
	},
}

func init() {
	spaceCmd.AddCommand(spaceLsCmd)
}
