package download

import "errors"

var (
	ErrEmpty        = errors.New("empty input")
	ErrUnrecognized = errors.New("unrecognized download input; use http(s) URL, magnet link, or .torrent path")
)
