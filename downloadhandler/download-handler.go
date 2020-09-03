package downloadhandler

import (
	"context"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
)

var (
	progressOutputs map[string]chan int64
	downloaders     map[string]*downloader.Downloader
)

func GetOrCreateDownloader(ctx context.Context, client client.Client, videoID string) (*downloader.Downloader, chan int64) {
	d, ok := downloaders[videoID]
	if !ok {
		newDownloader := downloader.New(client, videoID)
		downloaders[videoID] = newDownloader

		newProgressOutput := make(chan int64)
		progressOutputs[videoID] = newProgressOutput

		go startDownload(ctx, newDownloader, newProgressOutput)

		return newDownloader, newProgressOutput
	}

	p := progressOutputs[videoID]

	return d, p
}

func RemoveDownloader(videoID string) {
	if _, ok := downloaders[videoID]; ok {
		delete(downloaders, videoID)
	}

	if _, ok := progressOutputs[videoID]; ok {
		delete(progressOutputs, videoID)
	}
}

func startDownload(ctx context.Context, d *downloader.Downloader, p chan int64) {
	d.RetrieveVideoInfo(ctx)
	go d.DownloadVideo(ctx, p)
}
