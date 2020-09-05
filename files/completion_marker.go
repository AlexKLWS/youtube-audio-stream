package files

import (
	"os"
	"path/filepath"

	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/spf13/viper"
)

// CreateCompletionMarker creates an empty "completion marker" file in the target directory specified by path
// This file signifies that download/transmuxing has been complete
func CreateCompletionMarker(path string) error {
	f := filepath.Join(path, viper.GetString(consts.CompletionMarker))
	_, err := os.Create(f)
	return err
}
