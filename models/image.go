package models

type Image struct {
	UUID         string   `json:"uuid"`
	URLImage     string   `json:"url_image"`
	URLThumbnail string   `json:"url_thumbnail"`
	Hashtags     []string `json:"hashtags"`
	CreatedAt    string   `json:"xata_createdat"`
	SizeKb       int      `json:"size_kb"`
}
