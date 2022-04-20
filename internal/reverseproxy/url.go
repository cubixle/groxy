package reverseproxy

import (
	"net/url"
	"strings"
)

func buildReqURL(remoteAddr string, u *url.URL) (*url.URL, error) {
	newURL, err := url.Parse("//" + remoteAddr)
	if err != nil {
		return nil, err
	}

	newURL.Scheme = "http"
	newURL.RawQuery = u.RawQuery

	newURL.Path = strings.Join([]string{
		strings.TrimRight(newURL.Path, "/"),
		strings.TrimLeft(u.Path, "/"),
	}, "/")

	return newURL, nil
}
