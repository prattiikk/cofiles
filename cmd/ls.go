/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	auth "github.com/prattiikk/cofiles/cmd/utils/auth"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		header, err := auth.GetAuthHeader()
		if err != nil || header == "" {
			fmt.Println("Authorization failed:", err)
			return
		}

		config := auth.LoadConfig()
		url := config.Server + "/core/files"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Failed to create request:", err)
			return
		}

		req.Header.Set("Authorization", header)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("HTTP request failed:", err)
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response body:", err)
			return
		}

		// Struct to match the JSON response
		type File struct {
			Name      string `json:"name"`
			CreatedAt string `json:"createdAt"`
		}
		var files []File

		err = json.Unmarshal(bodyBytes, &files)
		if err != nil {
			fmt.Println("Failed to parse JSON response:", err)
			fmt.Println("Raw response:", string(bodyBytes))
			return
		}

		if len(files) == 0 {
			fmt.Println("No personal files found.")
			return
		}

		fmt.Println("Your Files:")
		for _, file := range files {
			fmt.Printf("- %s (Created at: %s)\n", file.Name, file.CreatedAt)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
