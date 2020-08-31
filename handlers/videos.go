package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
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
		log.Fatal(err)
		return err
	}
	defer ws.Close()

	_, url, err := ws.ReadMessage()
	if err != nil {
		log.Fatal(err)
		return err
	}

	if viper.GetBool(consts.Debug) {
		fmt.Printf("Downloading URL: %s\n", url)
	}

	queue := make(chan int64)

	outputURL, err := utils.FormOutputURL(string(url))
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err := ws.WriteJSON(models.ProgressUpdate{Type: models.DOWNLOAD_BEGUN, OutputURL: outputURL}); err != nil {
		log.Fatal(err)
		return err
	}
	c := client.Get()
	d := downloader.New(c, string(url))
	d.RetrieveVideoInfo(ctx.Request().Context())
	go d.DownloadVideo(ctx.Request().Context(), queue)
	for elem := range queue {
		if err := ws.WriteJSON(models.ProgressUpdate{Type: models.DOWNLOAD_IN_PROGRESS, DownloadPercentage: int(elem)}); err != nil {
			log.Fatal(err)
			return err
		}
	}
	outputDir := d.GetVideoID()
	sourceFilePath := d.GetVideoFilePath()

	if err := ws.WriteJSON(models.ProgressUpdate{Type: models.TRANSMUXING_BEGUN}); err != nil {
		log.Fatal(err)
		return err
	}

	t := transmuxer.New(outputDir, sourceFilePath)
	t.ConvertVideo()
	if err := ws.WriteJSON(models.ProgressUpdate{Type: models.TRANSMUXING_FINISHED}); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
