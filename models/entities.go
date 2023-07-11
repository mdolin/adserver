package models

// AdFormat represents the format of the ad
type AdFormat string

const (
	Banner       AdFormat = "banner"
	Interstitial AdFormat = "interstitial"
	Video        AdFormat = "video"
)

// AdUnit represents an advertising space
type AdUnit struct {
	ID     string
	Format AdFormat
	Width  int
	Height int
}

// Creative represents an ad that can be displayed
type Creative struct {
	ID      string
	Format  AdFormat
	Width   int
	Height  int
	Content string
	Price   float64
}

// AdRequest represents the structure of the incoming ad request JSON
type AdRequest struct {
	AdUnitID string `json:"ad_unit_id"`
	UserID   string `json:"user_id"`
}

// AdResponse represents the response format for the ad server
type AdResponse struct {
	CreativeID string  `json:"creative_id"`
	Content    string  `json:"content"`
	Price      float64 `json:"price"`
	UserID     string  `json:"user_id"`
}
