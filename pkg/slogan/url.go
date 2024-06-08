package slogan

import (
	"log/slog"
	"net/url"
)

func SanitizedURL(k string, u *url.URL) slog.Attr {
	// lazy staff, may be cur of user creds from url by another way
	cu := *u
	cu.User = nil

	return slog.String(k, cu.String())
}
