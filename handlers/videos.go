package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/downloader"
	"github.com/AlexKLWS/youtube-audio-stream/transmuxer"
	"github.com/labstack/echo"
)

type VideoDownloadRequest struct {
	URL string `json:"url" xml:"url"`
}

func DownloadAndProcessVideo(ctx echo.Context) error {
	requestBody := VideoDownloadRequest{}

	defer ctx.Request().Body.Close()
	err := json.NewDecoder(ctx.Request().Body).Decode(&requestBody)
	if err != nil {
		log.Printf("Failed processing video download request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	c := client.Get()
	d := downloader.New(c, requestBody.URL)
	d.DownloadVideo(ctx.Request().Context())
	outputDir := d.GetVideoID()
	sourceFilePath := d.GetVideoFilePath()

	t := transmuxer.New(outputDir, sourceFilePath)
	t.ConvertVideo()
	return ctx.String(http.StatusOK, "OK")
}
