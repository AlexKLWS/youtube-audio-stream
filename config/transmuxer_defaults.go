package config

// TransmuxerDefaults is default arguments and flags for ffmpeg file conversion
var TransmuxerDefaults = []string{
	"-c:a", "aac",
	"-b:a", "128k",
	"-map", "0:0",
	"-f", "hls",
	"-hls_list_size", "0",
	"-hls_time", "10",
	"-hls_segment_filename"}
