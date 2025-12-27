package models

// Platform represents a social media platform
type Platform string

const (
	PlatformTikTok    Platform = "tiktok"
	PlatformX         Platform = "x"
	PlatformInstagram Platform = "instagram"
	PlatformYouTube   Platform = "youtube"
)

// String returns the string representation of the platform
func (p Platform) String() string {
	return string(p)
}

// IsValid checks if the platform is valid
func (p Platform) IsValid() bool {
	switch p {
	case PlatformTikTok, PlatformX, PlatformInstagram, PlatformYouTube:
		return true
	default:
		return false
	}
}
