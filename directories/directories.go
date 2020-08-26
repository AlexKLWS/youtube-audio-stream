package directories

import (
	"log"
	"os"

	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/spf13/viper"
)

func PrepareDirectories() {
	if _, err := os.Stat(viper.GetString(consts.OutputDir)); os.IsNotExist(err) {
		if err2 := os.Mkdir(viper.GetString(consts.OutputDir), os.ModePerm); err2 != nil {
			log.Fatal(err2)
		}
	}

	if _, err := os.Stat(viper.GetString(consts.SourceDir)); os.IsNotExist(err) {
		if err2 := os.Mkdir(viper.GetString(consts.SourceDir), os.ModePerm); err2 != nil {
			log.Fatal(err2)
		}
	}
}
