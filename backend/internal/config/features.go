package config

// Feature Flag store
type Features struct {

	// feature flag on /llm/ask
	Agent bool `env:"AGENT_ENABLED, default=false"`
}
