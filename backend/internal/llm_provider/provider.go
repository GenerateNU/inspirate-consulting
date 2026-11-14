package llmprovider

type LLM string

const (
	GeminiModel LLM = "gemini"
)

type LLMProvider struct {
	UserInput string
	Model     LLM
}
