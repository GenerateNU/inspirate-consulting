package models

// the info needed to presign a profile picture upload for the logged-in user
type PresignPfpUploadRequestBody struct {
	ContentType   string `json:"content_type" enum:"image/jpeg,image/png,image/webp" doc:"MIME type of the image; the upload must send this exact Content-Type header"`
	ContentLength int64  `json:"content_length" minimum:"1" maximum:"5242880" doc:"Size of the image in bytes, at most 5 MiB; the upload must send exactly this many bytes"`
}

type PresignPfpUploadInput struct {
	Body PresignPfpUploadRequestBody
}

type PresignPfpUploadOutput struct {
	Body PresignUploadResponse
}

type UpdatePfpRequestBody struct {
	S3Key string `json:"s3_key" minLength:"1" doc:"The s3_key returned by PUT /user/pfp, after the image was uploaded to its upload_url"`
}

type UpdatePfpInput struct {
	Body UpdatePfpRequestBody
}

// the logged-in user's current profile picture
type Pfp struct {
	S3Key       string `json:"s3_key" doc:"The profile picture's S3 object key"`
	DownloadURL string `json:"download_url" doc:"A presigned, time-limited URL for the image. Re-fetch GET /user/me if the URL expires."`
}

type UpdatePfpOutput struct {
	Body Pfp
}
