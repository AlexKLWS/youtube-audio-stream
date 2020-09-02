package transmuxer

import (
	"log"
	"os"
	"os/exec"

	"github.com/AlexKLWS/youtube-audio-stream/config"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/spf13/viper"
	"golang.org/x/exp/errors/fmt"
)

type Transmuxer struct {
	arguments      []string
	sourceFilePath string
	outputDir      string
}

func New(outputDir string, sourceFilePath string) *Transmuxer {
	return &Transmuxer{outputDir: outputDir, sourceFilePath: sourceFilePath}
}

func (t *Transmuxer) ConvertVideo() error {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", viper.GetString(consts.OutputDir), t.outputDir)); os.IsNotExist(err) {
		if err2 := os.Mkdir(fmt.Sprintf("%s/%s", viper.GetString(consts.OutputDir), t.outputDir), os.ModePerm); err2 != nil {
			log.Print(err2)
			return err2
		}
	}

	args := append([]string{"-i", t.sourceFilePath}, config.TransmuxerDefaults...)
	segmentOutputFilename := fmt.Sprintf("%s/%s/out%%03d.ts", viper.GetString(consts.OutputDir), t.outputDir)
	args = append(args, segmentOutputFilename)
	playlistOutputFilename := fmt.Sprintf("%s/%s/out.m3u8", viper.GetString(consts.OutputDir), t.outputDir)
	args = append(args, playlistOutputFilename)

	cmd := exec.Command("ffmpeg", args...)
	err := cmd.Run()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
