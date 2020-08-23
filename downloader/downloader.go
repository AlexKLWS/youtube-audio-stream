package downloader

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/decipher"
	"github.com/AlexKLWS/youtube-audio-stream/exerrors"
	"github.com/AlexKLWS/youtube-audio-stream/models"
	videoinfo "github.com/AlexKLWS/youtube-audio-stream/video_info"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

// Downloader offers high level functions to download videos into files
type Downloader struct {
	client    *client.Client
	OutputDir string // optional directory to store the files
}

// New creates a new downloader with provided client
func New(c *client.Client) *Downloader {
	return &Downloader{client: c}
}

func (dl *Downloader) getOutputFile(v *videoinfo.VideoInfo, format *models.Format, outputFile string) (string, error) {

	if outputFile == "" {
		outputFile = SanitizeFilename(v.Title)
		outputFile += pickIdealFileExtension(format.MimeType)
	}

	if dl.OutputDir != "" {
		if err := os.MkdirAll(dl.OutputDir, 0755); err != nil {
			return "", err
		}
		outputFile = filepath.Join(dl.OutputDir, outputFile)
	}

	return outputFile, nil
}

//Download : Starting download video by arguments.
func (dl *Downloader) Download(ctx context.Context, v *videoinfo.VideoInfo, format *models.Format, outputFile string) error {
	destFile, err := dl.getOutputFile(v, format, outputFile)
	if err != nil {
		return err
	}

	// Create output file
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	dl.logf("Download to file=%s", destFile)
	return dl.videoDLWorker(ctx, out, v, format)
}

func (dl *Downloader) videoDLWorker(ctx context.Context, out *os.File, video *videoinfo.VideoInfo, format *models.Format) error {
	resp, err := dl.getStream(ctx, video, format)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	prog := &progress{
		contentLength: float64(resp.ContentLength),
	}

	// create progress bar
	progress := mpb.New(mpb.WithWidth(64))
	bar := progress.AddBar(
		int64(prog.contentLength),

		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)

	reader := bar.ProxyReader(resp.Body)
	mw := io.MultiWriter(out, prog)
	_, err = io.Copy(mw, reader)
	if err != nil {
		return err
	}

	progress.Wait()
	return nil
}

// GetStreamContext returns the HTTP response for a specific format with a context
func (dl *Downloader) getStream(ctx context.Context, video *videoinfo.VideoInfo, format *models.Format) (*http.Response, error) {
	url, err := dl.getStreamURL(ctx, video, format)
	if err != nil {
		return nil, err
	}

	return dl.client.HTTPGet(ctx, url)
}

func (dl *Downloader) getStreamURL(ctx context.Context, video *videoinfo.VideoInfo, format *models.Format) (string, error) {
	if format.URL != "" {
		return format.URL, nil
	}

	cipher := format.Cipher
	if cipher == "" {
		return "", exerrors.ErrCipherNotFound
	}

	return decipher.FormURLFromCipher(ctx, dl.client, video.ID, cipher)
}

func (dl *Downloader) logf(format string, v ...interface{}) {
	if !dl.client.Silent {
		log.Printf(format, v...)
	}
}
