package v1

import (
	"net/url"
)

func DeepCopyUrl(u *url.URL) *url.URL {
	result := *u
	copiedUser := *u.User
	result.User = &copiedUser
	return &result
}
