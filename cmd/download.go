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
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
