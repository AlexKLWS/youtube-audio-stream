package downloadhandler

import (
	"context"
	"sync"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
)

var (
	progressOutputs map[string]chan int64
	downloaders     map[string]*downloader.Downloader
	mutex           sync.RWMutex
)

func GetOrCreateDownloadHandle(ctx context.Context, client client.Client, videoID string) chan int64 {
	mutex.Lock()
	defer mutex.Unlock()

	if downloaders == nil {
		downloaders = make(map[string]*downloader.Downloader)
	}
	if progressOutputs == nil {
		progressOutputs = make(map[string]chan int64)
	}

	_, ok := downloaders[videoID]
	if !ok {
		newDownloader := downloader.New(client, videoID)
		downloaders[videoID] = newDownloader

		newProgressOutput := make(chan int64)
		progressOutputs[videoID] = newProgressOutput

		go startDownload(ctx, newDownloader, newProgressOutput)

		return newProgressOutput
	}

	p := progressOutputs[videoID]

	return p
}

func RemoveDownloader(videoID string) {
	mutex.Lock()
	defer mutex.Unlock()

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
