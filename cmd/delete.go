/*
Copyright © 2025 Pratik Shinde
*/
package cmd

import (
	"fmt"

	"github.com/prattiikk/cofiles/internal/cloud"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [filename]",
	Short: "Delete a file from your cloud storage",
	Long: `Delete a file from your personal cloud storage by specifying its filename.

This command will:
- Fetch your uploaded file list.
- Match the specified filename.
- If a match is found, it will be permanently deleted from both cloud storage and the database.

Usage:
  cofiles delete mynotes.txt

If the file is not found, you'll see a list of available files to choose from.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("❌ Please specify the filename to delete.")
			fmt.Println("Usage: cofiles delete <filename>")
			return
		}

		filename := args[0]

		files, err := cloud.GetUserFiles()
		if err != nil {
			fmt.Println("❌ Failed to fetch file list:", err)
			return
		}

		if len(files) == 0 {
			fmt.Println("📁 No files found in your cloud storage.")
			return
		}

		var matchedFile *cloud.File
		for _, f := range files {
			if f.Name == filename {
				matchedFile = &f
				break
			}
		}

		if matchedFile == nil {
			fmt.Printf("❌ File not found: %s\n", filename)
			fmt.Println("\nAvailable files:")
			for _, f := range files {
				fmt.Printf("  - %s (%.2f KB)\n", f.Name, float64(f.Size)/1024)
			}
			return
		}

		if err := cloud.DeleteFile(matchedFile.ID); err != nil {
			fmt.Println("❌ Failed to delete the file:", err)
			return
		}

		fmt.Printf("✅ %s (%.2f KB) deleted successfully!\n", matchedFile.Name, float64(matchedFile.Size)/1024)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}