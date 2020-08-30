package utils

import (
	"fmt"

	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/spf13/viper"
)

// FormOutputURL creates a url for the final output file
func FormOutputURL(url string) (string, error) {
	id, err := ExtractVideoID(url)
	if err != nil {
		return "", err
	}

	outputURL := fmt.Sprintf("%s/%s/out.m3u8", viper.GetString(consts.OutputRoute), id)

	return outputURL, nil
}
