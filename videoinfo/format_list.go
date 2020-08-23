package videoinfo

import "github.com/AlexKLWS/youtube-audio-stream/models"

type FormatList []models.Format

func (list FormatList) FindByItag(itagNo int) *models.Format {
	for i := range list {
		if list[i].ItagNo == itagNo {
			return &list[i]
		}
	}
	return nil
}
