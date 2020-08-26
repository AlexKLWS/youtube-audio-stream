package transmuxer

import (
	"log"
	"os"
	"os/exec"

	"github.com/AlexKLWS/youtube-audio-stream/config"
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

func (t *Transmuxer) ConvertVideo() {
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		if err2 := os.Mkdir("output", os.ModePerm); err2 != nil {
			log.Fatal(err2)
		}
	}
	if _, err := os.Stat(fmt.Sprintf("output/%s", t.outputDir)); os.IsNotExist(err) {
		if err2 := os.Mkdir(fmt.Sprintf("output/%s", t.outputDir), os.ModePerm); err2 != nil {
			log.Fatal(err2)
		}
	}
	lol := append([]string{"-i", t.sourceFilePath}, config.TransmuxerDefaults...)
	olo := fmt.Sprintf("output/%s/out%%03d.ts", t.outputDir)
	lol = append(lol, olo)
	olo = fmt.Sprintf("output/%s/out.m3u8", t.outputDir)
	lol = append(lol, olo)
	cmd := exec.Command("ffmpeg", lol...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
