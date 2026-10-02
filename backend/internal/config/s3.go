package config

// S3 holds S3-specific configuration
type S3 struct {
	Bucket string `env:"S3_BUCKET_NAME,default=inspirate-consulting-videos-478867930449-us-east-1-an"`
}
