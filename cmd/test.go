/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"

	auth "github.com/prattiikk/cofiles/cmd/utils/auth"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")

		header, err := auth.GetAuthHeader()
		if err != nil || header == "" {
			fmt.Println("error reading header")
		} else {
			fmt.Println("header:", header)
		}

		config := auth.LoadConfig()
		fmt.Println("config:")
		fmt.Println("JWT:", config.JWT)
		fmt.Println("Server:", config.Server)

		// Make HTTP GET request
		url := config.Server + "/core/protected" // Example: call /health endpoint
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Failed to create request:", err)
			return
		}

		// Attach auth header if available
		if header != "" {
			req.Header.Set("Authorization", header)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("HTTP request failed:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response body:", err)
			return
		}

		fmt.Println("Response status:", resp.Status)
		fmt.Println("Response body:", string(body))
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
