package videoinfo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"github.com/AlexKLWS/youtube-audio-stream/client"
	"github.com/AlexKLWS/youtube-audio-stream/exerrors"
	"github.com/AlexKLWS/youtube-audio-stream/models"
	"github.com/AlexKLWS/youtube-audio-stream/utils"
)

type VideoInfo struct {
	ID                  string
	Title               string
	Author              string
	Duration            time.Duration
	AllAvailableFormats FormatList
	Format              *models.Format
	FileFormat          string
}

// Fetch fetches video info metadata with a context
func Fetch(ctx context.Context, c client.Client, id string) (*VideoInfo, error) {
	// Circumvent age restriction to pretend access through googleapis.com
	eurl := "https://youtube.googleapis.com/v/" + id
	resp, err := c.HTTPGet(ctx, "https://youtube.com/get_video_info?video_id="+id+"&eurl="+eurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	v, err := parseVideoInfoString(string(body))
	if err != nil {
		return nil, err
	}

	v.ID = id

	return v, nil
}

// SelectFormat iterates over all available video formats and sets Format variable the most suitable one
func (v *VideoInfo) SelectFormat() {
	var audioOnlyFormats FormatList
	for _, f := range v.AllAvailableFormats {
		if checkIfAudio(f) {
			audioOnlyFormats = append(audioOnlyFormats, f)
		}
	}
	if len(audioOnlyFormats) != 0 {
		if bestFormat := audioOnlyFormats.FindByItag(140); bestFormat != nil {
			v.Format = bestFormat
			v.FileFormat = utils.PickIdealFileExtension(bestFormat.MimeType)
		}
	}
}

func parseVideoInfoString(info string) (v *VideoInfo, err error) {
	answer, err := url.ParseQuery(info)
	if err != nil {
		return nil, err
	}

	status := answer.Get("status")
	if status != "ok" {
		return nil, &exerrors.ErrResponseStatus{
			Status: status,
			Reason: answer.Get("reason"),
		}
	}

	// read the streams map
	playerResponse := answer.Get("player_response")
	if playerResponse == "" {
		return nil, errors.New("no player_response found in the server's answer")
	}

	var prData models.PlayerResponseData
	if err := json.Unmarshal([]byte(playerResponse), &prData); err != nil {
		return nil, fmt.Errorf("unable to parse player response JSON: %w", err)
	}

	v = &VideoInfo{}

	v.Title = prData.VideoDetails.Title
	v.Author = prData.VideoDetails.Author

	if seconds, _ := strconv.Atoi(prData.Microformat.PlayerMicroformatRenderer.LengthSeconds); seconds > 0 {
		v.Duration = time.Duration(seconds) * time.Second
	}

	// Check if video is downloadable
	if prData.PlayabilityStatus.Status != "OK" {
		return nil, &exerrors.ErrPlayabiltyStatus{
			Status: prData.PlayabilityStatus.Status,
			Reason: prData.PlayabilityStatus.Reason,
		}
	}

	// Assign Streams
	v.AllAvailableFormats = append(prData.StreamingData.Formats, prData.StreamingData.AdaptiveFormats...)

	if len(v.AllAvailableFormats) == 0 {
		return nil, errors.New("no formats found in the server's answer")
	}

	return
}

func checkIfAudio(f models.Format) bool {
	return f.FPS == 0 && f.Width == 0 && f.Height == 0
}
