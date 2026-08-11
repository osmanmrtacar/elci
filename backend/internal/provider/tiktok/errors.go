package tiktok

import (
	"errors"
	"fmt"
)

// friendlyPublishError rewrites the handful of TikTok error codes that mean
// "the account has hit a posting cap right now" into a message that tells
// the user to try again later, per TikTok's Direct Post guideline requiring
// clients to stop the attempt and communicate that plainly. Every other
// code (e.g. access_token_invalid) passes through unchanged — those are
// genuine errors, not something rewording would help.
func friendlyPublishError(err error) error {
	var apiErr apiError
	if !errors.As(err, &apiErr) {
		return err
	}

	switch apiErr.Code {
	case "spam_risk_too_many_posts", "reached_active_user_cap", "rate_limit_exceeded":
		return fmt.Errorf("TikTok has reached its posting limit for now — try again later (%s)", apiErr.Code)
	default:
		return err
	}
}
