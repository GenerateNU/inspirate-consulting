package llmprovider

import (
	"google.golang.org/genai"
	"google.golang.org/genai/interactions/models/interactions"
)

const (
	SearchTasksTool = "search_tasks"
	SearchChatsTool = "search_chats"
)

var tools = []interactions.Tool{
	interactions.NewTool(interactions.Function{
		Name:        genai.Ptr(SearchTasksTool),
		Description: genai.Ptr("Search a student's tasks by meaning. Returns matching tasks with due dates and completion status. Use for questions about deadlines, progress, or what is overdue."),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "What to search for in the students' tasks",
				},
			},
			"required": []string{"query"},
		},
	}),
	interactions.NewTool(interactions.Function{
		Name:        genai.Ptr(SearchChatsTool),
		Description: genai.Ptr("Search chat messages between a student and their counselor by meaning. Returns matching messages with sender and timestamp. Use for questions about what was discussed or said."),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "What to search for in the student and counselor chats",
				},
			},
			"required": []string{"query"},
		},
	}),
}
