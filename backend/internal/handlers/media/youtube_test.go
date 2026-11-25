package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseYoutubeVideoID(t *testing.T) {
	t.Parallel()

	valid := map[string]string{
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ":            "dQw4w9WgXcQ",
		"https://youtube.com/watch?v=dQw4w9WgXcQ&t=42s":          "dQw4w9WgXcQ",
		"http://m.youtube.com/watch?feature=share&v=dQw4w9WgXcQ": "dQw4w9WgXcQ",
		"https://youtu.be/dQw4w9WgXcQ":                           "dQw4w9WgXcQ",
		"https://youtu.be/dQw4w9WgXcQ?si=abc123":                 "dQw4w9WgXcQ",
		"  https://YouTu.be/dQw4w9WgXcQ  ":                       "dQw4w9WgXcQ",
	}
	for input, want := range valid {
		t.Run("valid "+input, func(t *testing.T) {
			t.Parallel()
			got, err := parseYoutubeVideoID(input)
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}

	invalid := []string{
		"",
		"dQw4w9WgXcQ",                          // bare id, no URL
		"javascript:alert(1)",                  // no host
		"https://evil.com/watch?v=dQw4w9WgXcQ", // wrong host
		"https://youtube.com.evil.com/watch?v=dQw4w9WgXcQ", // lookalike host
		"https://youtube.com@evil.com/watch?v=dQw4w9WgXcQ", // userinfo trick: the host is evil.com
		"https://www.youtube.com/watch",                    // missing v
		"https://www.youtube.com/watch?v=short",            // id wrong length
		"https://www.youtube.com/watch?v=dQw4w9WgXc%22",    // id with disallowed character
		"https://www.youtube.com/channel/UCabcdefghij",     // not a video
		"https://www.youtube.com/shorts/dQw4w9WgXcQ",       // only normal videos are supported
		"https://youtu.be/dQw4w9WgXcQ/extra",               // extra path segment
	}
	for _, input := range invalid {
		t.Run("invalid "+input, func(t *testing.T) {
			t.Parallel()
			_, err := parseYoutubeVideoID(input)
			assert.Error(t, err)
		})
	}
}
