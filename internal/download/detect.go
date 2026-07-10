package download

import (
	"os"
	"path/filepath"
	"strings"
)

type Kind int

const (
	KindHTTP Kind = iota
	KindMagnet
	KindTorrentFile
)

func (k Kind) String() string {
	switch k {
	case KindHTTP:
		return "http"
	case KindMagnet:
		return "magnet"
	case KindTorrentFile:
		return "torrent"
	default:
		return "unknown"
	}
}

type Input struct {
	Kind    Kind
	Raw     string
	Torrent string // local path when KindTorrentFile
}

func Parse(raw string) (Input, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return Input{}, ErrEmpty
	}
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "magnet:") {
		return Input{Kind: KindMagnet, Raw: s}, nil
	}
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return Input{Kind: KindHTTP, Raw: s}, nil
	}
	if strings.HasSuffix(lower, ".torrent") {
		path := expandHome(s)
		if _, err := os.Stat(path); err != nil {
			return Input{}, err
		}
		return Input{Kind: KindTorrentFile, Raw: s, Torrent: path}, nil
	}
	// Treat as path that might be a torrent file
	path := expandHome(s)
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return Input{Kind: KindTorrentFile, Raw: s, Torrent: path}, nil
	}
	return Input{}, ErrUnrecognized
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
