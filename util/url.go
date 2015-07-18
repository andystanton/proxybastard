package util

import (
	"fmt"
	"strings"
)

// SanitiseHTTPProxyURL strips the protocol (assumed to actually be http) and
// replaces it with http://, unless the string length without protocol is 0
// in which case the empty string is returned.
func SanitiseHTTPProxyURL(url string) string {
	noProtocol := url
	noProtocol = strings.TrimPrefix(noProtocol, "https://")
	noProtocol = strings.TrimPrefix(noProtocol, "http://")
	if len(noProtocol) > 0 {
		return fmt.Sprintf("http://%s", noProtocol)
	}
	return noProtocol
}
