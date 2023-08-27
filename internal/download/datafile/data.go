package datafile

type Data struct {
	URL      string                  `json:"url"`
	Site     string                  `json:"site"`
	Title    string                  `json:"title"`
	Type     DataType                `json:"type"`
	Streams  map[string]*Stream      `json:"streams"`
	Captions map[string]*CaptionPart `json:"caption"`
	Err      error                   `json:"err"`
}

type DataType string

const (
	DataTypeVideo DataType = "video"
	DataTypeImage DataType = "image"
	DataTypeAudio DataType = "audio"
)

type Stream struct {
	ID      string  `json:"id"`
	Quality string  `json:"quality"`
	Parts   []*Part `json:"parts"`
	Size    int64   `json:"size"`
	Ext     string  `json:"ext"`
	NeedMux bool    `json:"need_mux"`
}

type CaptionPart struct {
	Part
	Transform func([]byte) ([]byte, error) `json:"-"`
}

type Part struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
	Ext  string `json:"ext"`
}

// FillUpStreamsData fills up some data automatically.
func (d *Data) FillUpStreamsData() {
	for id, stream := range d.Streams {
		// fill up ID
		stream.ID = id
		if stream.Quality == "" {
			stream.Quality = id
		}

		// generate the merged file extension
		if d.Type == DataTypeVideo && stream.Ext == "" {
			ext := stream.Parts[0].Ext
			// The file extension in `Parts` is used as the merged file extension by default, except for the following formats
			switch ext {
			// ts and flv files should be merged into an mp4 file
			case "ts", "flv", "f4v":
				ext = "mp4"
			}
			stream.Ext = ext
		}

		// calculate total size
		if stream.Size > 0 {
			continue
		}
		var size int64
		for _, part := range stream.Parts {
			size += part.Size
		}
		stream.Size = size
	}
}

// EmptyData returns an "empty" Data object with the given URL and error.
func EmptyData(url string, err error) *Data {
	return &Data{
		URL: url,
		Err: err,
	}
}

// Options defines optional options that can be used in the extraction function.
type Options struct {
	// Playlist indicates if we need to extract the whole playlist rather than the single video.
	Playlist bool
	// Items defines wanted items from a playlist. Separated by commas like: 1,5,6,8-10.
	Items string
	// ItemStart defines the starting item of a playlist.
	ItemStart int
	// ItemEnd defines the ending item of a playlist.
	ItemEnd int

	// ThreadNumber defines how many threads will use in the extraction, only works when Playlist is true.
	ThreadNumber int
	Cookie       string

	// EpisodeTitleOnly indicates file name of each bilibili episode doesn't include the playlist title
	EpisodeTitleOnly bool

	YoukuCcode    string
	YoukuCkey     string
	YoukuPassword string
}

// Extractor implements video data extraction related operations.
type Extractor interface {
	// Extract is the main function to extract the data.
	Extract(url string, option Options) ([]*Data, error)
}
