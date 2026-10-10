package config

import "inspirate-consulting/internal/supabase"

// Config holds the entire application configuration
type Config struct {
	Application Application
	DB          DB
	Supabase    supabase.SupabaseInterface
	S3          S3
	// TestMode=true will skip the auth middleware
	TestMode bool
}
