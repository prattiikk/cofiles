/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/prattiikk/cofiles/internal/cloud"
	"github.com/spf13/cobra"
)

type File struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mimeType"`
	CreatedAt string `json:"createdAt"`
}

type FileListResponse struct {
	Success bool   `json:"success"`
	Files   []File `json:"files"`
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <filename>",
	Short: "Download a specific file from your cloud storage",
	Long: `Download a file stored in your authenticated cloud storage.

You need to specify the exact filename of the file you want to download. 
If the file is found, it will be downloaded using its ID. 
In case the file is not found, a list of all available files will be displayed.

Usage:
  cofiles download <filename>

Example:
  cofiles download notes.txt
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("❌ Please specify the filename to download.")
			fmt.Println("Usage: cofiles clouddown <filename>")
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

		fmt.Printf("✅ %s (%.2f KB) downloaded succssfully!\n", matchedFile.Name, float64(matchedFile.Size)/1024)

		if err := cloud.DownloadFile(matchedFile.ID); err != nil {
			fmt.Println("❌ Download failed:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
