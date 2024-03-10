package helpers

import (
	"fmt"
	"krzysztofjanowski/camera-surveillance-dashboard/packages/videos"
	"strings"
	"time"
)

// finds matching videos based on start and end dates, used for dasbhoard filtering
func FindMatchingVideos(ffromDate string, ftoDate string, videoDir string ) ([]string, error) {
	// slice container for videos that will match from and end dates 
	var filteredVideos []string 

	// format of the date, necessary for time.Parse 
	layout := "2006-01-02"
	 
	// dates coming from the dashboard page  converted into time object 
	startDateOnForm, _:= time.Parse(layout, ffromDate)
	endDateOnForm, _ := time.Parse(layout,  ftoDate)

	all_videos,err  := videos.ListVideos(videoDir)

	if err != nil {
		// fatal error importing videos
		return []string{}, err
	}

	for _, filename := range(all_videos){
		// extract date from video file names which are in format of YYYY-MM-DD|HH-MM-SS e.g. 2021-01-10|16:32:55.mp4
		date := strings.Split(filename, "|")[0]

		// Parse the input string into a time.Time object
		videoDate, err := time.Parse(layout, date)

		if err != nil {
			// Handle error if parsing fails
			// fmt.Println("Error parsing date:", err)
			return []string{}, err
		} else {

			if videoDate.After(startDateOnForm) && videoDate.Before(endDateOnForm){
				filteredVideos = append(filteredVideos, filename )
			}

		}
	}

	// print names of videos that match the start and end date 
	if len(filteredVideos) > 0 {
		fmt.Println("Following videos match the dates specified:")

		for _, filename := range(filteredVideos){
			fmt.Println("filename:", filename)
		}	
	} else {
		// "No videos with matching dates found"
	}

	return filteredVideos, nil 

}