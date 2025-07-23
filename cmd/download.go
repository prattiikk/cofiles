/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/prattiikk/cofiles/internal/auth"
	"github.com/prattiikk/cofiles/internal/storage"
	"github.com/spf13/cobra"
)

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

		// get the jwt header from config file
		header, err := auth.GetAuthHeader()
		if err != nil || header == "" {
			fmt.Println("Authorization failed:", err)
			return
		}

		// get the backend server url from the config
		config := auth.LoadConfig()

		// generate a endpoint to hit
		url := config.Server + "/core/files"

		// create a request object of type GET with the built url
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Failed to create request:", err)
			return
		}

		// attach the jwt header to the reqest object
		req.Header.Set("Authorization", header)

		// create a http client
		client := &http.Client{}

		// make the reqest and get the response
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("HTTP request failed:", err)
			return
		}

		defer resp.Body.Close()

		// read the response body
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("Failed to read response body:", err)
			return
		}

		// Structs to match the JSON response
		type File struct {
			Name      string `json:"name"`
			CreatedAt string `json:"createdAt"`
		}

		type FileListResponse struct {
			Files []File `json:"files"`
		}

		var fileList FileListResponse

		// parse the body which is in json and store it in fileList obejct
		err = json.Unmarshal(bodyBytes, &fileList)
		if err != nil {
			fmt.Println("Failed to parse JSON response:", err)
			fmt.Println("Raw response:", string(bodyBytes))
			return
		}

		if len(fileList.Files) == 0 {
			fmt.Println("No personal files found.")
			return
		}

		// // fmt.Println("Your Files:")
		// // print the list of files
		// for _, file := range fileList.Files {
		// 	fmt.Printf(" - %s\t - %s\n", file.Name, file.CreatedAt)
		// }

		filename := args[0]
		fmt.Printf("downloading %s", filename)
		errfile := storage.DownloadFile(filename)
		if errfile != nil {
			fmt.Println("❌ Error:", err)
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
