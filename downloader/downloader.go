package downloader

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/decipher"
	"github.com/AlexKLWS/youtube-audio-stream/exerrors"
	"github.com/AlexKLWS/youtube-audio-stream/files"
	"github.com/AlexKLWS/youtube-audio-stream/utils"
	"github.com/AlexKLWS/youtube-audio-stream/videoinfo"
	"github.com/reactivex/rxgo"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

// Downloader offers high level functions to download videos into files
type Downloader struct {
	videoID         string
	video           *videoinfo.VideoInfo
	client          client.Client
	outputDirectory string
	progressOutput  chan rxgo.Item
}

// New creates a new downloader with provided client
func New(c client.Client, videoID string, progressOutput chan rxgo.Item) *Downloader {
	return &Downloader{client: c, videoID: videoID, progressOutput: progressOutput}
}

// RetrieveVideoInfo fetches video info from youtube API
func (dl *Downloader) RetrieveVideoInfo(ctx context.Context) error {
	v, err := videoinfo.Fetch(ctx, dl.client, dl.videoID)
	if err != nil {
		return err
	}

	dl.video = v
	v.SelectFormat()

	return nil
}

//DownloadVideo returns a download handle
func (dl *Downloader) DownloadVideo(ctx context.Context) error {
	file, err := dl.getFileHandle()
	if err != nil {
		return err
	}

	defer file.Close()

	return dl.videoDLWorker(ctx, file)
}

func (dl *Downloader) getFileHandle() (*os.File, error) {
	destFile, err := dl.getOutputFilePath()
	if err != nil {
		return nil, err
	}

	// Create output file
	file, err := os.Create(destFile)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (dl *Downloader) getOutputFilePath() (string, error) {
	outputFilePath := utils.SanitizeFilename(dl.video.Title)
	outputFilePath += dl.video.FileFormat

	dl.outputDirectory = filepath.Join(viper.GetString(consts.SourceDir), dl.video.ID)

	if _, err := os.Stat(dl.outputDirectory); os.IsNotExist(err) {
		if err2 := os.Mkdir(dl.outputDirectory, os.ModePerm); err2 != nil {
			log.Print(err2)
			return "", err2
		}
	}
	outputFilePath = filepath.Join(dl.outputDirectory, outputFilePath)

	return outputFilePath, nil
}

func (dl *Downloader) videoDLWorker(ctx context.Context, file *os.File) error {
	resp, err := dl.getStream(ctx)
	if err != nil {
		log.Print(err)
		return err
	}
	defer resp.Body.Close()

	var src io.Reader

	// Send download data updates to progress output channel if it's available
	if dl.progressOutput != nil {
		writeCounter := &DownloadProgressWriter{ProgressOutput: dl.progressOutput}
		writeCounter.ContentLength = resp.ContentLength

		src = io.TeeReader(resp.Body, writeCounter)

		if _, err = io.Copy(file, src); err != nil {
			log.Print(err)
			return err
		}
	} else { // Otherwise print out progress to terminal
		prog := &ProgressBarWriter{
			contentLength: float64(resp.ContentLength),
		}
		progress := mpb.New(mpb.WithWidth(64))
		bar := progress.AddBar(
			int64(prog.contentLength),

			mpb.BarStyle("╢▌▌░╟"),
			mpb.PrependDecorators(
				decor.CountersKibiByte("% .2f / % .2f"),
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.EwmaETA(decor.ET_STYLE_GO, 90),
				decor.Name(" | "),
				decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
			),
		)

		reader := bar.ProxyReader(resp.Body)
		mw := io.MultiWriter(file, prog)
		if _, err = io.Copy(mw, reader); err != nil {
			log.Print(err)
			return err
		}
	}

	err = files.CreateCompletionMarker(dl.outputDirectory)

	return err
}

// GetStreamContext returns the HTTP response for a specific format with a context
func (dl *Downloader) getStream(ctx context.Context) (*http.Response, error) {
	url, err := dl.getStreamURL(ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return dl.client.HTTPGet(ctx, url)
}

func (dl *Downloader) getStreamURL(ctx context.Context) (string, error) {
	if dl.video.Format.URL != "" {
		return dl.video.Format.URL, nil
	}

	cipher := dl.video.Format.Cipher
	if cipher == "" {
		return "", exerrors.ErrCipherNotFound
	}

	return decipher.FormURLFromCipher(ctx, dl.client, dl.video.ID, cipher)
}
