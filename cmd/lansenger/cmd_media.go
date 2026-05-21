package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "Upload and download media files",
}

var mediaUploadCmd = &cobra.Command{
	Use:   "upload FILE_PATH",
	Short: "Upload a media file",
	Args:  cobra.ExactArgs(1),
	Run:   runMediaUpload,
}

var mediaDownloadCmd = &cobra.Command{
	Use:   "download MEDIA_ID",
	Short: "Download media by ID (outputs raw data info)",
	Args:  cobra.ExactArgs(1),
	Run:   runMediaDownload,
}

var mediaDownloadToFileCmd = &cobra.Command{
	Use:   "download-to-file MEDIA_ID",
	Short: "Download media and save to a file",
	Args:  cobra.ExactArgs(1),
	Run:   runMediaDownloadToFile,
}

var (
	mediaUploadType      int
	mediaDownloadOutput  string
	mediaDownloadTypeStr string
)

func init() {
	mediaUploadCmd.Flags().IntVarP(&mediaUploadType, "media-type", "t", 3, "Media type: 1=video, 2=image, 3=file")

	mediaDownloadToFileCmd.Flags().StringVarP(&mediaDownloadOutput, "output", "o", "", "Target file path (defaults to media ID)")
	mediaDownloadToFileCmd.Flags().StringVar(&mediaDownloadTypeStr, "media-type", "file", "Media type hint: file, image, video")

	mediaCmd.AddCommand(mediaUploadCmd)
	mediaCmd.AddCommand(mediaDownloadCmd)
	mediaCmd.AddCommand(mediaDownloadToFileCmd)
	rootCmd.AddCommand(mediaCmd)
}

func runMediaUpload(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.UploadMedia(ctx, args[0], mediaUploadType)
	checkError(err)
	outputResultFields(result, []string{"media_id"})
}

func runMediaDownload(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.DownloadMedia(ctx, args[0])
	checkError(err)

	if jsonOutput {
		m := map[string]interface{}{
			"success": result.Success,
			"size":    len(result.Data),
			"error":   result.Error,
		}
		outputJSON(m)
		return
	}

	if result.Success {
		fmt.Printf("Downloaded %d bytes\n", len(result.Data))
	} else {
		fmt.Printf("Error: %s\n", result.Error)
	}
}

func runMediaDownloadToFile(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	savedPath, err := client.DownloadMediaToFile(ctx, args[0], mediaDownloadOutput)
	checkError(err)

	if jsonOutput {
		m := map[string]interface{}{
			"saved_path":  savedPath,
			"media_type":  mediaDownloadTypeStr,
		}
		outputJSON(m)
		return
	}

	fmt.Printf("Saved to: %s\n", savedPath)
}

func mediaTypeFromString(s string) int {
	switch s {
	case "image":
		return 2
	case "video":
		return 1
	case "file":
		return 3
	default:
		n, err := strconv.Atoi(s)
		if err == nil {
			return n
		}
		return 3
	}
}