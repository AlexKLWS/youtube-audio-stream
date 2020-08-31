package models

// ProgressUpdateType is a enum describing how vide is being processed
type ProgressUpdateType uint

const (
	DOWNLOAD_BEGUN ProgressUpdateType = iota
	DOWNLOAD_IN_PROGRESS
	DOWNLOAD_FINISHED
	TRANSMUXING_BEGUN
	TRANSMUXING_FINISHED
)

// ProgressUpdate is a json object describing current state of video processing progress
type ProgressUpdate struct {
	Type               ProgressUpdateType `json:"type"`
	OutputURL          string             `json:"outputURL"`
	DownloadPercentage int                `json:"downloadPercentage"`
}
