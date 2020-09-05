package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/spf13/viper"
)

func GetSourceFilePath(videoID string) (string, error) {
	path := filepath.Join(viper.GetString(consts.SourceDir), videoID)
	var file string
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if info.Name() != viper.GetString(consts.CompletionMarker) && !info.IsDir() {
			// We assume there're only 2 files in a folder
			file = info.Name()
		}
		return nil
	})
	if file == "" {
		return "", fmt.Errorf("Actually no source file for id: %s", videoID)
	}
	path = filepath.Join(path, file)
	return path, err
}
