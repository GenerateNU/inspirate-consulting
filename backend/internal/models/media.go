package models

type Media struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	LengthInMins int `json:"length_in_mins"`
	SchoolYear *int `json:"school_year"`
	S3Key string `json:"s3_key"`
}

type MediaAccess struct {
	ID string `json:"id"`
	StudentID string `json:"student_id"`
	MediaID string `json:"media_id"`
}

type CreateMediaRequestBody struct{
	Title string `json:"title"`
	Description string `json:"description"`
	LengthInMins int `json:"length_in_mins"`
	SchoolYear *int `json:"school_year"`
	S3Key string `json:"s3_key"`
}

type GrantMediaAccessRequestBody struct{
	StudentID string `json:"student_id"`
	MediaID string `json:"media_id"`
}

type CreateMediaInput struct{
	Body CreateMediaRequestBody
}

type CreateMediaOutput struct{
	Body Media
}

type GetMediaInput struct{
	ID   string `path:"id"`
}

type GetMediaOutput struct{
	Body Media
}

type ListAllMediaInput struct{

}

type ListAllMediaOutput struct{
	Body []Media
}

type DeleteMediaInput struct{
	ID   string `path:"id"`
}


type GrantMediaAccessInput struct {
	Body GrantMediaAccessRequestBody
}

type GrantMediaAccessOutput struct {
	Body MediaAccess
}

type RevokeMediaAccessInput struct{
	ID   string `path:"id"`
}


type ListAccessibleMediaInput struct{

}

type ListAccessibleMediaOutput struct {
	Body []Media
}