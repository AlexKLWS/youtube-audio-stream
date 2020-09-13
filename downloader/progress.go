package downloader

import (
	"github.com/AlexKLWS/youtube-audio-stream/feed"
	"github.com/AlexKLWS/youtube-audio-stream/models"
)

// DownloadProgressWriter writes download percentage to provided ProgressOutput
type DownloadProgressWriter struct {
	totalDownloaded    int64
	previousPercentage int64
	ContentLength      int64
	ProgressOutput     *feed.ProgressUpdateFeed
}

// Write implements the io.Writer interface.
//
// Always completes and never returns an error.
func (wc *DownloadProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	wc.totalDownloaded += int64(n)
	percentage := wc.totalDownloaded * 100 / wc.ContentLength
	if percentage != wc.previousPercentage {
		wc.ProgressOutput.Send(models.ProgressUpdate{Type: models.DOWNLOAD_IN_PROGRESS, DownloadPercentage: int(percentage)})
		wc.previousPercentage = percentage
	}
	return n, nil
}

// ProgressBarWriter is used to write download progress to mpb progress bar
type ProgressBarWriter struct {
	contentLength     float64
	totalWrittenBytes float64
	downloadLevel     float64
}

func (pbc *ProgressBarWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	pbc.totalWrittenBytes = pbc.totalWrittenBytes + float64(n)
	currentPercent := (pbc.totalWrittenBytes / pbc.contentLength) * 100
	if (pbc.downloadLevel <= currentPercent) && (pbc.downloadLevel < 100) {
		pbc.downloadLevel++
	}
	return
}
