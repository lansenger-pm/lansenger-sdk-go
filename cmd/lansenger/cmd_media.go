package main

import (
	"context"
	"fmt"
	"os"
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

var mediaPathCmd = &cobra.Command{
	Use:   "path MEDIA_ID",
	Short: "Fetch media file download path info",
	Args:  cobra.ExactArgs(1),
	Run:   runMediaPath,
}

var mediaUploadAppCmd = &cobra.Command{
	Use:   "upload-app FILE_PATH",
	Short: "Upload app/bot media (4.5.4 endpoint)",
	Args:  cobra.ExactArgs(1),
	Run:   runMediaUploadApp,
}

var mediaUploadAppV2Cmd = &cobra.Command{
	Use:   "upload-app-v2 FILE_PATH",
	Short: "Upload app/bot media V2 (4.5.5 endpoint, requires user token)",
	Args:  cobra.ExactArgs(1),
	Run:   runMediaUploadAppV2,
}

var mediaDownloadShareCmd = &cobra.Command{
	Use:   "download-share SHARE_ID",
	Short: "Download media by share ID (4.5.6)",
	Args:  cobra.ExactArgs(1),
	Run:   runMediaDownloadShare,
}

var (
	mediaUploadType         int
	mediaUploadUserToken    string
	mediaDownloadOutput     string
	mediaDownloadTypeStr    string
	mediaPathUserToken string

	mediaUploadAppType     string
	mediaUploadAppWidth    int
	mediaUploadAppHeight   int
	mediaUploadAppDuration int
)

var (
	mediaUploadAppV2Type      string
	mediaUploadAppV2UserToken string
	mediaUploadAppV2Width     int
	mediaUploadAppV2Height    int
	mediaUploadAppV2Duration  int

	mediaDownloadShareOutput    string
	mediaDownloadShareUserToken string
)

func init() {
	mediaUploadCmd.Flags().IntVarP(&mediaUploadType, "media-type", "t", 2, "Media type: 1=video, 2=image, 3=audio")
	mediaUploadCmd.Flags().StringVar(&mediaUploadUserToken, "user-token", "", "User token")

	mediaDownloadToFileCmd.Flags().StringVarP(&mediaDownloadOutput, "output", "o", "", "Target file path (defaults to media ID)")
	mediaDownloadToFileCmd.Flags().StringVar(&mediaDownloadTypeStr, "media-type", "file", "Media type hint: file, image, video")

	mediaPathCmd.Flags().StringVar(&mediaPathUserToken, "user-token", "", "User token")

	mediaUploadAppCmd.Flags().StringVarP(&mediaUploadAppType, "media-type", "t", "file", "Media type: file, video, image, audio")
	mediaUploadAppCmd.Flags().IntVar(&mediaUploadAppWidth, "width", 0, "Width (for video/image)")
	mediaUploadAppCmd.Flags().IntVar(&mediaUploadAppHeight, "height", 0, "Height (for video/image)")
	mediaUploadAppCmd.Flags().IntVar(&mediaUploadAppDuration, "duration", 0, "Duration in seconds (for video/audio)")

	mediaUploadAppV2Cmd.Flags().StringVarP(&mediaUploadAppV2Type, "media-type", "t", "file", "Media type: file, video, image, audio")
	mediaUploadAppV2Cmd.Flags().StringVar(&mediaUploadAppV2UserToken, "user-token", "", "User token (required)")
	mediaUploadAppV2Cmd.Flags().IntVar(&mediaUploadAppV2Width, "width", 0, "Width (for video/image)")
	mediaUploadAppV2Cmd.Flags().IntVar(&mediaUploadAppV2Height, "height", 0, "Height (for video/image)")
	mediaUploadAppV2Cmd.Flags().IntVar(&mediaUploadAppV2Duration, "duration", 0, "Duration in seconds (for video/audio)")

	mediaDownloadShareCmd.Flags().StringVarP(&mediaDownloadShareOutput, "output", "o", "", "Target file path")
	mediaDownloadShareCmd.Flags().StringVar(&mediaDownloadShareUserToken, "user-token", "", "User token")

	mediaCmd.AddCommand(mediaUploadCmd)
	mediaCmd.AddCommand(mediaDownloadCmd)
	mediaCmd.AddCommand(mediaDownloadToFileCmd)
	mediaCmd.AddCommand(mediaPathCmd)
	mediaCmd.AddCommand(mediaUploadAppCmd)
	mediaCmd.AddCommand(mediaUploadAppV2Cmd)
	mediaCmd.AddCommand(mediaDownloadShareCmd)
	rootCmd.AddCommand(mediaCmd)
}

func runMediaUpload(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.UploadMedia(ctx, args[0], mediaUploadType, mediaUploadUserToken)
	checkError(err)
	outputResultFields(result, []string{"media_id", "created_time"})
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

func runMediaPath(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchMediaPath(ctx, args[0], mediaPathUserToken)
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

func runMediaUploadAppV2(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.UploadAppMediaV2(ctx, args[0], mediaUploadAppV2Type, mediaUploadAppV2UserToken, mediaUploadAppV2Width, mediaUploadAppV2Height, mediaUploadAppV2Duration)
	checkError(err)
	outputResultFields(result, []string{"media_id"})
}

func runMediaDownloadShare(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.DownloadMediaByShareID(ctx, args[0], mediaDownloadShareUserToken)
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
		if mediaDownloadShareOutput != "" {
			err := os.WriteFile(mediaDownloadShareOutput, result.Data, 0644)
			checkError(err)
			fmt.Printf("Saved to: %s (%d bytes)\n", mediaDownloadShareOutput, len(result.Data))
		} else {
			fmt.Printf("Downloaded %d bytes\n", len(result.Data))
		}
	} else {
		fmt.Printf("Error: %s\n", result.Error)
	}
}