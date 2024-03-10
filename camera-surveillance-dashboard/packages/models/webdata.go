package models

// bucket for data sent to webpages (from handlers to template)
// struct to allow variety of data types be passed to template package 
type WebData struct {
	Notifications	[]string // for notification messages  
	LatestVideo 	string 
	VideosSlice		[]string
	MovementSensor1 bool
	MovementSensor2 bool
	LightSensor		int
	CrossSiteToken  string
	StringMap		map[string]string //  
}
