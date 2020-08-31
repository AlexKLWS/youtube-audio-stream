package models

type ProgressUpdateType int

const (
	DOWNLOAD_BEGUN ProgressUpdateType = iota
	DOWNLOAD_IN_PROGRESS
	DOWNLOAD_FINISHED
	TRANSMUXING_BEGUN
	TRANSMUXING_FINISHED
)

type ProgressUpdate struct {
	Type               ProgressUpdateType `json:"type"`
	OutputURL          string             `json:"outputURL"`
	DownloadPercentage int                `json:"downloadPercentage"`
}
