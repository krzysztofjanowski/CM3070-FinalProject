package videos 


import (
	"os"
	"path/filepath"
)

// list all the videos in the video dir
func ListVideos(videoDir string) ([]string, error) {

	// videoDir = "../../web/ready_videos"
	var videos []string
 
	// the directory that contains the video
	videoDirectory, error := os.Open(videoDir)

	if error != nil {
		return nil, error
	}

	defer videoDirectory.Close()

	files, error := videoDirectory.ReadDir(-1)

	if error != nil {
		return nil, error
	}

	for _, file := range files {
		//check if file is mp4
		if filepath.Ext(file.Name()) == ".mp4" {
			videos = append(videos, file.Name())
		}
	}

	return videos, nil

}
