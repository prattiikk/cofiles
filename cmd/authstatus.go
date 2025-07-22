/*
Copyright © 2025 NAME
*/
package cmd

import (
	"fmt"

	auth "github.com/prattiikk/cofiles/cmd/utils"
	"github.com/spf13/cobra"
)

// authStatusCmd represents the authStatus command
var authStatusCmd = &cobra.Command{
	Use:   "authstatus",
	Short: "Check authentication status",
	Long:  `Displays whether the user is authenticated and prints the stored JWT if available.`,
	Run: func(cmd *cobra.Command, args []string) {
		if auth.IsAuthenticated() {
			token, err := auth.GetJWT()
			if err != nil || token == "" {
				fmt.Println("You are authenticated, but no valid token was found.")
				fmt.Println("run : cofile login")
				return
			}
			fmt.Println("✅ Authenticated")
			fmt.Println("JWT Token:", token)
		} else {
			fmt.Println("❌ Not authenticated")
		}
	},
}

func init() {
	rootCmd.AddCommand(authStatusCmd)
}
