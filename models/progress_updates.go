package models

// ProgressUpdateType is a enum describing how vide is being processed
type ProgressUpdateType uint

const (
	ERROR ProgressUpdateType = iota
	REQUEST_ACCEPTED
	DOWNLOAD_BEGUN
	DOWNLOAD_IN_PROGRESS
	DOWNLOAD_FINISHED
	TRANSMUXING_BEGUN
	TRANSMUXING_FINISHED
	AUDIO_IS_AVAILABLE
)

// ProgressUpdate is a json object describing current state of video processing progress
type ProgressUpdate struct {
	Type               ProgressUpdateType `json:"type"`
	VideoID            string             `json:"videoID"`
	DownloadPercentage int                `json:"downloadPercentage"`
	Error              error              `json:"error"`
}
