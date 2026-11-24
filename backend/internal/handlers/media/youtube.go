package media

import (
	"net/url"
	"regexp"
	"strings"

	"inspirate-consulting/internal/errs"
)

// YouTube video ids are always 11 characters from this set. Because the frontend puts the id into an
// iframe src, guarantees only YouTube videos can be embedded.
var youtubeVideoIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

// extracts the video id from a normal YouTube link: the address bar link (youtube.com/watch?v=ID)
// or the Share button link (youtu.be/ID)
func parseYoutubeVideoID(rawURL string) (string, error) {
	invalid := errs.BadRequest("youtube_url must be a youtube.com/watch or youtu.be video link")

	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", invalid
	}

	var id string
	switch strings.ToLower(u.Hostname()) {
	case "youtu.be":
		id = strings.TrimPrefix(u.Path, "/")
	case "youtube.com", "www.youtube.com", "m.youtube.com":
		if u.Path == "/watch" {
			id = u.Query().Get("v")
		}
	}

	if !youtubeVideoIDPattern.MatchString(id) {
		return "", invalid
	}
	return id, nil
}
