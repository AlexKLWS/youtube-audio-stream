package downloader

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/decipher"
	"github.com/AlexKLWS/youtube-audio-stream/exerrors"
	"github.com/AlexKLWS/youtube-audio-stream/utils"
	"github.com/AlexKLWS/youtube-audio-stream/videoinfo"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

// Downloader offers high level functions to download videos into files
type Downloader struct {
	url       string
	video     *videoinfo.VideoInfo
	client    client.Client
	OutputDir string // optional directory to store the files
}

// New creates a new downloader with provided client
func New(c client.Client, url string) *Downloader {
	return &Downloader{client: c, url: url}
}

//DownloadVideo returns a download handle
func (dl *Downloader) DownloadVideo(ctx context.Context, outputFilename string) error {
	v, err := videoinfo.Fetch(ctx, dl.client, dl.url)
	if err != nil {
		return err
	}
	dl.video = v
	v.SelectFormat()

	destFile, err := dl.getOutputFile(outputFilename)
	if err != nil {
		return err
	}

	// Create output file
	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	return dl.videoDLWorker(ctx, out)
}

func (dl *Downloader) getOutputFile(outputFile string) (string, error) {

	if outputFile == "" {
		outputFile = utils.SanitizeFilename(dl.video.Title)
		outputFile += dl.video.FileFormat
	}

	if dl.OutputDir != "" {
		if err := os.MkdirAll(dl.OutputDir, 0755); err != nil {
			return "", err
		}
		outputFile = filepath.Join(dl.OutputDir, outputFile)
	}

	return outputFile, nil
}

func (dl *Downloader) videoDLWorker(ctx context.Context, out *os.File) error {
	resp, err := dl.getStream(ctx)
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

		mpb.BarStyle("╢▌▌░╟"),
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
func (dl *Downloader) getStream(ctx context.Context) (*http.Response, error) {
	url, err := dl.getStreamURL(ctx)
	if err != nil {
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
