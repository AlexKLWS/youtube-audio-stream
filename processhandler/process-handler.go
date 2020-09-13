package processhandler

import (
	"context"
	"log"
	"sync"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
	"github.com/AlexKLWS/youtube-audio-stream/files"
	"github.com/AlexKLWS/youtube-audio-stream/models"
	"github.com/AlexKLWS/youtube-audio-stream/transmuxer"
	"github.com/reactivex/rxgo"
	"github.com/spf13/viper"
)

var (
	progressOutputs map[string]rxgo.Observable
	mutex           sync.RWMutex
)

func GetOrCreateProcessHandle(client client.Client, videoID string) rxgo.Observable {
	mutex.Lock()
	defer mutex.Unlock()

	if progressOutputs == nil {
		progressOutputs = make(map[string]rxgo.Observable)
	}

	p, ok := progressOutputs[videoID]
	if !ok {
		newProgressOutput := make(chan rxgo.Item)
		progressOutputs[videoID] = rxgo.FromEventSource(newProgressOutput, rxgo.WithBackPressureStrategy(rxgo.Drop))

		newDownloader := downloader.New(client, videoID, newProgressOutput)

		go handleProcessing(videoID, newDownloader, newProgressOutput)

		return progressOutputs[videoID]
	}

	return p
}

func handleProcessing(videoID string, d *downloader.Downloader, p chan rxgo.Item) {
	defer removeProcessHandler(videoID)
	defer close(p)

	ctx := context.Background()

	if viper.GetBool(consts.Debug) {
		log.Printf("Downloading video with id: %s\n", videoID)
	}

	if !files.CheckIfWasProcessed(viper.GetString(consts.SourceDir), videoID) {
		p <- rxgo.Of(models.ProgressUpdate{Type: models.DOWNLOAD_BEGUN, VideoID: videoID})
		if err := d.RetrieveVideoInfo(ctx); err != nil {
			log.Print(err)
			p <- rxgo.Of(models.ProgressUpdate{Type: models.ERROR, Error: err})
			return
		}
		if err := d.DownloadVideo(ctx); err != nil {
			log.Print(err)
			p <- rxgo.Of(models.ProgressUpdate{Type: models.ERROR, Error: err})
			return
		}
		p <- rxgo.Of(models.ProgressUpdate{Type: models.DOWNLOAD_FINISHED})
	}

	p <- rxgo.Of(models.ProgressUpdate{Type: models.TRANSMUXING_BEGUN})
	sourceFilePath, err := files.GetSourceFilePath(videoID)
	if err != nil {
		log.Print(err)
		p <- rxgo.Of(models.ProgressUpdate{Type: models.ERROR, Error: err})
		return
	}
	t := transmuxer.New(videoID, sourceFilePath)
	if err := t.ConvertVideo(); err != nil {
		log.Print(err)
		p <- rxgo.Of(models.ProgressUpdate{Type: models.ERROR, Error: err})
		return
	}
	p <- rxgo.Of(models.ProgressUpdate{Type: models.TRANSMUXING_FINISHED})

	p <- rxgo.Of(models.ProgressUpdate{Type: models.AUDIO_IS_AVAILABLE, VideoID: videoID})

	if viper.GetBool(consts.Debug) {
		log.Printf("Video %s\n is ready for streaming!", videoID)
	}
}

func removeProcessHandler(videoID string) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := progressOutputs[videoID]; ok {
		delete(progressOutputs, videoID)
	}
}
