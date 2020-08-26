package main

import (
	"github.com/AlexKLWS/youtube-audio-stream/cmd"
	"github.com/AlexKLWS/youtube-audio-stream/config"
)

func main() {
	config.InitializeViper()

	cmd.SetupRootCommand()
	cmd.SetupServeCommand()
	cmd.Execute()
}
