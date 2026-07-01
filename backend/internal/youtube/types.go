package youtube

// Holds information about youtube video
type YTVideo struct {
	ID          string
	Title       string
	URL         string
	Description string
	Duration    float64
	Channel     string
	ViewCount   int
	Thumbnails  []VideoThumbnail
}

// Holds information about YouTube channel
type YTChannel struct {
	ID   string
	Name string
	URL  string
}

// Holds information about video thumbnail
type VideoThumbnail struct {
	URL    string
	Height int
	Width  int
}

// Used to unmarshal the JSON output from yt-dlp command
// JSON tags correspond to fields in yt-dlp JSON output
type YTDLPVideo struct {
	ID          string           	  `json:"id"`
	Title       string           	  `json:"title"`
	URL         string           	  `json:"webpage_url"`
	Description string           	  `json:"description"`
	Duration    float64          	  `json:"duration"`
	Channel     string           	  `json:"channel"`
	ViewCount   int             	  `json:"view_count"`
	Thumbnails  []YTDLPVideoThumbnail `json:"thumbnails"`
}

// Used to unmarshal the thumbnails field from yt-dlp JSON output
type YTDLPVideoThumbnail struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
