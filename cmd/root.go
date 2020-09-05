package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
	"github.com/AlexKLWS/youtube-audio-stream/files"
	"github.com/AlexKLWS/youtube-audio-stream/transmuxer"
	"github.com/AlexKLWS/youtube-audio-stream/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "youtube-audio-stream",
		Short: "Runs the server",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runRoot(args[0])
		},
	}
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}

// SetupRootCommand binds command line flags and arguments to viper config instance
func SetupRootCommand() {
	rootCmd.PersistentFlags().BoolP("debug", "d", viper.GetString("env") == "debug", "run in debug mode")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func runRoot(url string) {
	if viper.GetBool(consts.Debug) {
		fmt.Println("Running in debug mode")
	}
	files.PrepareDirectories()
	httpTransport := client.GetHTTPTransport()

	c := client.New(httpTransport)
	ctx := context.Background()

	videoID, err := utils.ExtractVideoID(string(url))
	if err != nil {
		log.Fatal(err)
	}

	d := downloader.New(c, videoID)
	d.RetrieveVideoInfo(ctx)
	d.DownloadVideo(ctx, nil)
	outputDir := videoID
	sourceFilePath := d.GetVideoFilePath()

	t := transmuxer.New(outputDir, sourceFilePath)
	t.ConvertVideo()
}
