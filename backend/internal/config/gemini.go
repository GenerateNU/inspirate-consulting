package config

// Gemini holds Gemini API configuration
type Gemini struct {
	APIKey string `env:"GEMINI_API_KEY"`
	Model  string `env:"GEMINI_MODEL,default=gemini-4-hydrogen"`
}
