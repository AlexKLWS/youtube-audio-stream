package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
	"github.com/AlexKLWS/youtube-audio-stream/transmuxer"
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

	c := client.Get()
	d := downloader.New(c, string(url))
	if err := ws.WriteMessage(websocket.TextMessage, []byte("Downloading video...")); err != nil {
		log.Fatal(err)
		return err
	}
	d.RetrieveVideoInfo(ctx.Request().Context())
	go d.DownloadVideo(ctx.Request().Context(), queue)
	for elem := range queue {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(strconv.FormatInt(elem, 10))); err != nil {
			log.Fatal(err)
			return err
		}
	}
	outputDir := d.GetVideoID()
	sourceFilePath := d.GetVideoFilePath()

	if err := ws.WriteMessage(websocket.TextMessage, []byte("Converting video...")); err != nil {
		log.Fatal(err)
		return err
	}

	t := transmuxer.New(outputDir, sourceFilePath)
	t.ConvertVideo()
	if err := ws.WriteMessage(websocket.TextMessage, []byte("Video converted successfully...")); err != nil {
		log.Fatal(err)
		return err
	}

	ws.Close()
	return nil
}
