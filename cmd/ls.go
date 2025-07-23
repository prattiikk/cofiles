/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/prattiikk/cofiles/internal/cloud"
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

		files, err := cloud.GetUserFiles()
		if err != nil {
			fmt.Println("❌ Failed to fetch file list:", err)
			return
		}

		if len(files) == 0 {
			fmt.Println("📁 No files found in your cloud storage.")
			return
		}

		// Print the files
		fmt.Printf("\n📂 Your Files:\n\n")
		fmt.Printf("%-30s %-20s %-10s %-25s\n", "Name", "Type", "Size (KB)", "Created At")
		fmt.Println("--------------------------------------------------------------------------------------")
		for _, file := range files {
			fmt.Printf("%-20s %-30s %-10.2f %-30s\n",
				file.Name, file.MimeType, float64(file.Size)/1024, file.CreatedAt)
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
