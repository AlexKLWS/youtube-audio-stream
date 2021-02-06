package processhandler

import (
	"context"
	"log"
	"sync"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
	"github.com/AlexKLWS/youtube-audio-stream/feed"
	"github.com/AlexKLWS/youtube-audio-stream/files"
	"github.com/AlexKLWS/youtube-audio-stream/models"
	"github.com/AlexKLWS/youtube-audio-stream/transmuxer"
	"github.com/spf13/viper"
)

var (
	feeds map[string]*feed.ProgressUpdateFeed
	mutex sync.RWMutex
)

func GetOrCreateSubscription(client client.Client, videoID string) <-chan models.ProgressUpdate {
	mutex.Lock()
	defer mutex.Unlock()

	if feeds == nil {
		feeds = make(map[string]*feed.ProgressUpdateFeed)
	}

	p, ok := feeds[videoID]
	if !ok {
		feed := feed.New()
		feeds[videoID] = feed

		newDownloader := downloader.New(client, videoID, feed)

		go handleProcessing(videoID, newDownloader, feed)

		return feed.Subscribe()
	}

	return p.Subscribe()
}

func handleProcessing(videoID string, d *downloader.Downloader, f *feed.ProgressUpdateFeed) {
	defer removeFeed(videoID)
	defer f.Close()

	ctx := context.Background()

	if viper.GetBool(consts.Debug) {
		log.Printf("Downloading video with id: %s\n", videoID)
	}

	if !files.CheckIfWasProcessed(viper.GetString(consts.SourceDir), videoID) {
		f.Send(models.ProgressUpdate{Type: models.DOWNLOAD_BEGUN, VideoID: videoID})
		if err := d.RetrieveVideoInfo(ctx); err != nil {
			log.Print(err)
			f.Send(models.ProgressUpdate{Type: models.ERROR, Error: err})
			return
		}
		if err := d.DownloadVideo(ctx); err != nil {
			log.Print(err)
			f.Send(models.ProgressUpdate{Type: models.ERROR, Error: err})
			return
		}
		f.Send(models.ProgressUpdate{Type: models.DOWNLOAD_FINISHED})
	}

	f.Send(models.ProgressUpdate{Type: models.TRANSMUXING_BEGUN})
	sourceFilePath, err := files.GetSourceFilePath(videoID)
	if err != nil {
		log.Print(err)
		f.Send(models.ProgressUpdate{Type: models.ERROR, Error: err})
		return
	}
	t := transmuxer.New(videoID, sourceFilePath)
	if err := t.ConvertVideo(); err != nil {
		log.Print(err)
		f.Send(models.ProgressUpdate{Type: models.ERROR, Error: err})
		return
	}
	f.Send(models.ProgressUpdate{Type: models.TRANSMUXING_FINISHED})

	f.Send(models.ProgressUpdate{Type: models.AUDIO_IS_AVAILABLE, VideoID: videoID})

	if viper.GetBool(consts.Debug) {
		log.Printf("Video %s\n is ready for streaming!", videoID)
	}
}

func removeFeed(videoID string) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := feeds[videoID]; ok {
		delete(feeds, videoID)
	}
}
