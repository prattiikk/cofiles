/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	auth "github.com/prattiikk/cofiles/internal/auth"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Sign out from your Cofiles account",
	Long: `Clears your local authentication session and signs you out of the Cofiles CLI.

After running this command, you will need to log in again using 'cofiles login' to perform any authenticated operations.

Example:
  cofiles logout
`,
	Run: func(cmd *cobra.Command, args []string) {
		auth.ClearAuth()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
