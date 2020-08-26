package handlers

import "github.com/AlexKLWS/youtube-audio-stream/router"

func RegisterHandlers(serverRouter *router.Router) {

	serverRouter.Videos.POST("", DownloadAndProcessVideo)

}
