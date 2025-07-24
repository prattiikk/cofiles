/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/prattiikk/cofiles/internal/cloud"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to your cloud storage",
	Long: `Upload a local file to your cloud storage.

Provide the full or relative path to the file you wish to upload.

Examples:
  cofile upload notes.txt
  cofile upload ./projects/report.pdf

Notes:
- Ensure you're authenticated before running this command.
- You’ll receive a confirmation message on successful upload.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide the path to the file to upload.")
			return
		}
		filePath := args[0]

		err := cloud.UploadFile(filePath)

		if err != nil {
			fmt.Printf("Error uploading file: %v\n", err)
			return
		}
		fmt.Println("✅ File uploaded successfully!")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
