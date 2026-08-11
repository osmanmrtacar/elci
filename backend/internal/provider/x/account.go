package x

import (
	"context"

	"github.com/osmanmertacar/elci/backend/internal/provider"
)

type userResponse struct {
	Data struct {
		ID              string `json:"id"`
		Username        string `json:"username"`
		Name            string `json:"name"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"data"`
}

func (p *Provider) AccountInfo(ctx context.Context, token provider.Token) (provider.AccountInfo, error) {
	var res userResponse
	err := bearerJSON(ctx, "GET", token.AccessToken, "/2/users/me?user.fields=profile_image_url", nil, &res)
	if err != nil {
		return provider.AccountInfo{}, err
	}

	return provider.AccountInfo{
		PlatformUserID: res.Data.ID,
		Username:       res.Data.Username,
		DisplayName:    res.Data.Name,
		AvatarURL:      res.Data.ProfileImageURL,
		Extra:          map[string]any{},
	}, nil
}
