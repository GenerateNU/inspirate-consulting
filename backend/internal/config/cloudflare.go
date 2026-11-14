package config

// Cloudflare holds Workers AI configuration used for text embeddings.
type Cloudflare struct {
	AccountID string `env:"CLOUDFLARE_ACCOUNT_ID"`
	APIToken  string `env:"CLOUDFLARE_API_TOKEN"`
}
