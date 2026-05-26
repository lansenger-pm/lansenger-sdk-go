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

var mediaFetchPathCmd = &cobra.Command{
	Use:   "fetch-path MEDIA_ID",
	Short: "Fetch media file download path info",
	Args:  cobra.ExactArgs(1),
	Run:   runMediaFetchPath,
}

var mediaUploadAppCmd = &cobra.Command{
	Use:   "upload-app FILE_PATH",
	Short: "Upload app/bot media (4.5.4 endpoint)",
	Args:  cobra.ExactArgs(1),
	Run:   runMediaUploadApp,
}

var (
	mediaUploadType         int
	mediaUploadUserToken    string
	mediaDownloadOutput     string
	mediaDownloadTypeStr    string
	mediaFetchPathUserToken string

	mediaUploadAppType     string
	mediaUploadAppWidth    int
	mediaUploadAppHeight   int
	mediaUploadAppDuration int
)

func init() {
	mediaUploadCmd.Flags().IntVarP(&mediaUploadType, "media-type", "t", 3, "Media type: 1=video, 2=image, 3=file, 4=audio")
	mediaUploadCmd.Flags().StringVar(&mediaUploadUserToken, "user-token", "", "User token")

	mediaDownloadToFileCmd.Flags().StringVarP(&mediaDownloadOutput, "output", "o", "", "Target file path (defaults to media ID)")
	mediaDownloadToFileCmd.Flags().StringVar(&mediaDownloadTypeStr, "media-type", "file", "Media type hint: file, image, video")

	mediaFetchPathCmd.Flags().StringVar(&mediaFetchPathUserToken, "user-token", "", "User token")

	mediaUploadAppCmd.Flags().StringVarP(&mediaUploadAppType, "media-type", "t", "file", "Media type: file, video, image, audio")
	mediaUploadAppCmd.Flags().IntVar(&mediaUploadAppWidth, "width", 0, "Width (for video/image)")
	mediaUploadAppCmd.Flags().IntVar(&mediaUploadAppHeight, "height", 0, "Height (for video/image)")
	mediaUploadAppCmd.Flags().IntVar(&mediaUploadAppDuration, "duration", 0, "Duration in seconds (for video/audio)")

	mediaCmd.AddCommand(mediaUploadCmd)
	mediaCmd.AddCommand(mediaDownloadCmd)
	mediaCmd.AddCommand(mediaDownloadToFileCmd)
	mediaCmd.AddCommand(mediaFetchPathCmd)
	mediaCmd.AddCommand(mediaUploadAppCmd)
	rootCmd.AddCommand(mediaCmd)
}

func runMediaUpload(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.UploadMedia(ctx, args[0], mediaUploadType, mediaUploadUserToken)
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

func runMediaFetchPath(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchMediaPath(ctx, args[0], mediaFetchPathUserToken)
	checkError(err)
	outputResultFields(result, []string{"media_path", "name", "type", "size"})
}

func runMediaUploadApp(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.UploadAppMedia(ctx, args[0], mediaUploadAppType, mediaUploadAppWidth, mediaUploadAppHeight, mediaUploadAppDuration)
	checkError(err)
	outputResultFields(result, []string{"media_id"})
}