package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/directories"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
	"github.com/AlexKLWS/youtube-audio-stream/transmuxer"
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
	rootCmd.PersistentFlags().BoolP("debug", "d", viper.GetString("env") == "debug", "specify server port")
	viper.BindPFlags(serveCmd.Flags())
}

func runRoot(url string) {
	directories.PrepareDirectories()
	httpTransport := client.GetHTTPTransport()

	c := client.New(httpTransport)
	ctx := context.Background()

	d := downloader.New(c, url)
	d.RetrieveVideoInfo(ctx)
	d.DownloadVideo(ctx, nil)
	outputDir := d.GetVideoID()
	sourceFilePath := d.GetVideoFilePath()

	t := transmuxer.New(outputDir, sourceFilePath)
	t.ConvertVideo()
}
