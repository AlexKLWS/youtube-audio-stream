package handlers

import (
	"log"
	"net/http"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/files"
	"github.com/AlexKLWS/youtube-audio-stream/models"
	"github.com/AlexKLWS/youtube-audio-stream/processhandler"
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

	videoID, err := utils.ExtractVideoID(string(url))
	if err != nil {
		log.Print(err)
		if err := ws.WriteJSON(models.ProgressUpdate{Type: models.ERROR}); err != nil {
			log.Print(err)
		}
		return nil
	}

	if files.CheckIfWasProcessed(viper.GetString(consts.OutputDir), videoID) {
		if ws.WriteJSON(models.ProgressUpdate{Type: models.AUDIO_IS_AVAILABLE, VideoID: videoID}); err != nil {
			log.Print(err)
		}
		return nil
	}

	if ws.WriteJSON(models.ProgressUpdate{Type: models.REQUEST_ACCEPTED, VideoID: videoID}); err != nil {
		log.Print(err)
		return nil
	}

	progressUpdates := processhandler.GetOrCreateProcessHandle(client.Get(), videoID)
	for update := range progressUpdates {
		if err := ws.WriteJSON(update); err != nil {
			log.Print(err)
			break
		}
	}

	return nil
}
