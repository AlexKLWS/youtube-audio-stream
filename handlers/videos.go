package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/downloadhandler"
	"github.com/AlexKLWS/youtube-audio-stream/files"
	"github.com/AlexKLWS/youtube-audio-stream/models"
	"github.com/AlexKLWS/youtube-audio-stream/transmuxer"
	"github.com/AlexKLWS/youtube-audio-stream/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

// DownloadAndProcessVideo performs video download and transmuxing, sending updates to client via sockets during the process
func DownloadAndProcessVideo(ctx echo.Context) error {
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer ws.Close()

	_, url, err := ws.ReadMessage()
	if err != nil {
		log.Print(err)
		if err := ws.WriteJSON(models.ProgressUpdate{Type: models.ERROR}); err != nil {
			log.Print(err)
		}
		return nil
	}

	if viper.GetBool(consts.Debug) {
		fmt.Printf("Downloading URL: %s\n", url)
	}

	videoID, err := utils.ExtractVideoID(string(url))
	if err != nil {
		log.Print(err)
		if err := ws.WriteJSON(models.ProgressUpdate{Type: models.ERROR}); err != nil {
			log.Print(err)
		}
		return nil
	}

	if files.CheckIfWasProcessed(viper.GetString(consts.OutputDir), videoID) {
		ws.WriteJSON(models.ProgressUpdate{Type: models.AUDIO_IS_ALREADY_AVAILABLE, VideoID: videoID})
		return nil
	}

	if !files.CheckIfWasProcessed(viper.GetString(consts.SourceDir), videoID) {
		if err := ws.WriteJSON(models.ProgressUpdate{Type: models.DOWNLOAD_BEGUN, VideoID: videoID}); err != nil {
			log.Print(err)
			return err
		}

		progressUpdates := downloadhandler.GetOrCreateDownloadHandle(ctx.Request().Context(), client.Get(), videoID)
		for update := range progressUpdates {
			if err := ws.WriteJSON(models.ProgressUpdate{Type: models.DOWNLOAD_IN_PROGRESS, DownloadPercentage: int(update)}); err != nil {
				log.Print(err)
			}
		}

		downloadhandler.RemoveDownloader(videoID)
	}

	if err := ws.WriteJSON(models.ProgressUpdate{Type: models.TRANSMUXING_BEGUN}); err != nil {
		log.Print(err)
	}

	sourceFilePath, err := files.GetSourceFilePath(videoID)
	if err != nil {
		log.Print(err)
		if err := ws.WriteJSON(models.ProgressUpdate{Type: models.ERROR}); err != nil {
			log.Print(err)
		}
		return nil
	}

	t := transmuxer.New(videoID, sourceFilePath)
	t.ConvertVideo()
	if err := ws.WriteJSON(models.ProgressUpdate{Type: models.TRANSMUXING_FINISHED}); err != nil {
		log.Print(err)
	}

	if viper.GetBool(consts.Debug) {
		fmt.Printf("Video %s\n is ready for streaming!", videoID)
	}

	return nil
}
