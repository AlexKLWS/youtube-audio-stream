package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
	videoinfo "github.com/AlexKLWS/youtube-audio-stream/video_info"
	"github.com/spf13/cobra"
)

// the command to run the server
var rootCmd = &cobra.Command{
	Use:   "youtube-audio-stream",
	Short: "Runs the server",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(args[0])
	},
}

func main() {
	// config.InitializeViper()
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
	v, err := videoinfo.Fetch(ctx, c, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
	}
	d := downloader.New(c)
	d.Download(ctx, v, &v.Formats[0], "lmao.mp4")
}
