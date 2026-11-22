package models

import "time"

type Year string

const (
	Freshman  Year = "freshman"
	Sophomore Year = "sophomore"
	Junior    Year = "junior"
	Senior    Year = "senior"
)


// Exactly one of S3Key / YoutubeVideoID is set (enforced by the media_single_source constraint).
// Field order must match the SELECT column order, since list queries use pgx.RowToStructByPos.
type Media struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	LengthInMins   *int    `json:"length_in_mins"`
	SchoolYear     *Year   `json:"school_year"`
	S3Key          *string `json:"s3_key"`
	YoutubeVideoID *string `json:"youtube_video_id" doc:"11-character YouTube video id; embed as https://www.youtube.com/embed/{id}"`
}

type MediaAccess struct {
	ID          string     `json:"id"`
	StudentID   string     `json:"student_id"`
	MediaID     string     `json:"media_id"`
	StartedAt   *time.Time `json:"started_at" doc:"When the student first opened the media"`
	CompletedAt *time.Time `json:"completed_at" doc:"When the student first marked the media complete"`
}

// AccessibleMedia is one row of a student's media list: the media plus that student's progress on it.
// Field order must match list_accessible_media.sql, since it is scanned with pgx.RowToStructByPos.
type AccessibleMedia struct {
	Media
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	DueDate     *time.Time `json:"due_date" doc:"Deadline of the student's todo that assigned this media (the most recent one, if assigned more than once)"`
}

type CreateMediaRequestBody struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	LengthInMins *int    `json:"length_in_mins,omitempty"`
	SchoolYear   *Year   `json:"school_year"`
	S3Key        *string `json:"s3_key,omitempty" minLength:"1" doc:"S3 key of an uploaded video; provide exactly one of s3_key or youtube_url"`
	YoutubeURL   *string `json:"youtube_url,omitempty" doc:"youtube.com or youtu.be video link; only the parsed video id is stored. Provide exactly one of s3_key or youtube_url"`
}

// UpdateMediaFields is what the store applies on PATCH /media/{id}; nil fields are left unchanged.
// At most one of S3Key / YoutubeVideoID is set, and setting one clears the other.
type UpdateMediaFields struct {
	Title          *string
	Description    *string
	S3Key          *string
	YoutubeVideoID *string
}

type GrantMediaAccessRequestBody struct {
	StudentID string `json:"student_id"`
	MediaID   string `json:"media_id"`
}

type CreateMediaInput struct {
	Body CreateMediaRequestBody
}

type CreateMediaOutput struct {
	Body Media
}

type UpdateMediaInputBody struct {
	Title       *string `json:"title,omitempty" minLength:"1"`
	Description *string `json:"description,omitempty"`
	S3Key       *string `json:"s3_key,omitempty" minLength:"1" doc:"S3 key of an uploaded video; replaces the video source (clears youtube_video_id). Mutually exclusive with youtube_url"`
	YoutubeURL  *string `json:"youtube_url,omitempty" doc:"youtube.com or youtu.be video link; replaces the video source (clears s3_key). Mutually exclusive with s3_key"`
}

type UpdateMediaInput struct {
	ID   string `path:"id"`
	Body UpdateMediaInputBody
}

type UpdateMediaOutput struct {
	Body Media
}

type GetMediaInput struct {
	ID string `path:"id"`
}

type GetMediaOutput struct {
	Body Media
}

type ListAllMediaInput struct {
	Query  string `query:"q" required:"false" doc:"Only return media whose title or description is similar to this text (pg_trgm similarity), best match first"`
	Limit  int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Maximum number of media to return"`
	Offset int    `query:"offset" default:"0" minimum:"0" doc:"Number of media to skip before returning results"`
}

type ListAllMediaOutput struct {
	Body []Media
}

type DeleteMediaInput struct {
	ID string `path:"id"`
}

type GrantMediaAccessInput struct {
	Body GrantMediaAccessRequestBody
}

type GrantMediaAccessOutput struct {
	Body MediaAccess
}

type RevokeMediaAccessInput struct {
	ID string `path:"id"`
}

type ListAccessibleMediaInput struct {
	Query  string `query:"q" required:"false" doc:"Only return media whose title or description is similar to this text (pg_trgm similarity), best match first"`
	Limit  int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Maximum number of media to return"`
	Offset int    `query:"offset" default:"0" minimum:"0" doc:"Number of media to skip before returning results"`
}

type ListAccessibleMediaOutput struct {
	Body []AccessibleMedia
}

type MediaProgressStatus string

const (
	MediaProgressStarted   MediaProgressStatus = "started"
	MediaProgressCompleted MediaProgressStatus = "completed"
)

type RecordMediaProgressInput struct {
	ID   string `path:"id"`
	Body struct {
		Status MediaProgressStatus `json:"status" enum:"started,completed" doc:"started when the student opens the media; completed when they mark it done. Repeating a status keeps the original timestamp"`
	}
}

type RecordMediaProgressOutput struct {
	Body MediaAccess
}
