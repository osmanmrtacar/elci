package tiktok

import (
	"context"
	"net/url"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

type userInfoResponse struct {
	User struct {
		OpenID      string `json:"open_id"`
		DisplayName string `json:"display_name"`
		AvatarURL   string `json:"avatar_url"`
	} `json:"user"`
}

type creatorInfoResponse struct {
	CreatorAvatarURL        string   `json:"creator_avatar_url"`
	CreatorUsername         string   `json:"creator_username"`
	CreatorNickname         string   `json:"creator_nickname"`
	PrivacyLevelOptions     []string `json:"privacy_level_options"`
	CommentDisabled         bool     `json:"comment_disabled"`
	DuetDisabled            bool     `json:"duet_disabled"`
	StitchDisabled          bool     `json:"stitch_disabled"`
	MaxVideoPostDurationSec int      `json:"max_video_post_duration_sec"`
}

// AccountInfo merges TikTok's basic profile (nickname/avatar) with a fresh
// creator_info query, since TikTok requires the Direct Post UX to show live
// privacy/interaction capabilities on every visit to the composer, not a
// cached copy from when the account was first connected.
func (p *Provider) AccountInfo(ctx context.Context, token provider.Token) (provider.AccountInfo, error) {
	var userInfo userInfoResponse
	query := url.Values{"fields": {"open_id,display_name,avatar_url"}}
	if err := getJSON(ctx, token.AccessToken, "/user/info/", query, &userInfo); err != nil {
		return provider.AccountInfo{}, err
	}

	var creatorInfo creatorInfoResponse
	if err := postJSON(ctx, token.AccessToken, "/post/publish/creator_info/query/", nil, &creatorInfo); err != nil {
		return provider.AccountInfo{}, err
	}

	return provider.AccountInfo{
		PlatformUserID: userInfo.User.OpenID,
		Username:       creatorInfo.CreatorUsername,
		// The Direct Post UX guideline requires showing the creator's
		// nickname specifically from creator_info on the posting page, not
		// whatever /user/info/ happens to have cached.
		DisplayName: creatorInfo.CreatorNickname,
		AvatarURL:   userInfo.User.AvatarURL,
		Extra: map[string]any{
			"privacy_level_options":       creatorInfo.PrivacyLevelOptions,
			"comment_disabled":            creatorInfo.CommentDisabled,
			"duet_disabled":               creatorInfo.DuetDisabled,
			"stitch_disabled":             creatorInfo.StitchDisabled,
			"max_video_post_duration_sec": creatorInfo.MaxVideoPostDurationSec,
		},
	}, nil
}
