package models

type AskLLMInput struct {
	Body struct {
		Input string `json:"input" minLength:"1" doc:"The question to ask the model"`
	}
}

type AskLLMOutput struct {
	Body struct {
		Answer string `json:"answer" doc:"The model's response"`
	}
}
