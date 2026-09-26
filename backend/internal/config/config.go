package config

// Config holds the entire application configuration
type Config struct {
	Application Application
	DB          DB
	Supabase    Supabase
	S3    		S3
	// TestMode=true will skip the auth middleware
	TestMode bool
}
