/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	auth "github.com/prattiikk/cofiles/internal/auth"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate your CLI with your cloud account",
	Long: `Log in to your Cofiles account using your browser.

This command opens a browser window with authentication options (e.g., Google or GitHub).
Once authenticated, your session token will be securely saved and used for future CLI interactions.

Example:
  cofiles login
`,
	Run: func(cmd *cobra.Command, args []string) {
		auth.Authenticate("http://localhost:3000")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
