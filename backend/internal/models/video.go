package models

// represents a single object in the video bucket
type Video struct {
	S3Key string `json:"s3_key" doc:"The video's S3 object key/path, e.g. video-bucket/common-app-tips.mp4"`
	DownloadURL string `json:"download_url" doc:"A presigned, time-limited URL for playback. Re-fetch this endpoint if the URL expires."`
}

// the info needed to generate a unique S3 key and presigned upload URL.
type PresignUploadRequestBody struct {
	OriginalFilename string `json:"original_filename" doc:"The filename of the video about to be uploaded"`
}

// what is returned after requesting an upload.
// separated from Video because the URL is for uploading (PUT), not downloading (GET)
type PresignUploadResponse struct {
	S3Key string `json:"s3_key" doc:"The S3 key/path to store elsewhere, e.g. in a media table"`
	UploadURL string `json:"upload_url" doc:"A presigned URL the client PUTs the raw video bytes to"`
}

// identifies a video by its S3 key. 
// taken as a query parameter rather than a path parameter to avoid ambiguity with path segment splitting.
type GetVideoInput struct {
	S3Key string `query:"s3_key" doc:"The S3 key/path of the video to fetch"`
}


// Huma readable input and output models for the upload, get, and list operations for videos in the S3 bucket.
type PresignUploadInput struct {
	Body PresignUploadRequestBody
}

type PresignUploadOutput struct {
	Body PresignUploadResponse
}

type GetVideoOutput struct {
	Body Video
}

// lists every object currently in the bucket, requires no input.
type ListVideosInput struct{}

type ListVideosOutput struct {
	Body []Video
}