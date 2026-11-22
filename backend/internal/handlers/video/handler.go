package video

import storage "inspirate-consulting/internal/data"

type Handler struct {
	VideoRepository       storage.VideoRepository
	MediaAccessRepository storage.MediaAccessRepository
}

// Here is where we create the handler which has reference to the VideoRepository
func NewHandler(videoRepository storage.VideoRepository, mediaAccessRepository storage.MediaAccessRepository) *Handler {
	return &Handler{
		VideoRepository:       videoRepository,
		MediaAccessRepository: mediaAccessRepository,
	}
}
