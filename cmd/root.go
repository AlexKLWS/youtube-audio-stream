package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
	"github.com/AlexKLWS/youtube-audio-stream/transmuxer"
	"github.com/spf13/cobra"
)

var (
	outputFilename string
	rootCmd        = &cobra.Command{
		Use:   "youtube-audio-stream",
		Short: "Runs the server",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			run(args[0])
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&outputFilename, "output", "o", "", "specify output filename")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}

func run(url string) {
	httpTransport := &http.Transport{
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpTransport.DialContext = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext

	c := client.New(httpTransport)
	ctx := context.Background()

	d := downloader.New(c, url)
	d.DownloadVideo(ctx)
	outputDir := d.GetVideoID()
	sourceFilePath := d.GetVideoFilePath()

	t := transmuxer.New(outputDir, sourceFilePath)
	t.ConvertVideo()
}
