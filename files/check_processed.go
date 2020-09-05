package files

import (
	"os"
	"path/filepath"

	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/spf13/viper"
)

// CheckIfWasProcessed checks if completion marker has been created
func CheckIfWasProcessed(outerFolder string, videoID string) bool {
	path := filepath.Join(outerFolder, videoID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	result := false
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if info.Name() == viper.GetString(consts.CompletionMarker) {
			result = true
		}
		return nil
	})
	return result
}
