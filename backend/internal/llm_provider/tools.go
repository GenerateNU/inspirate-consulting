package llmprovider

import (
	"google.golang.org/genai"
	"google.golang.org/genai/interactions/models/interactions"
)

const (
	SearchTasksTool = "search_tasks"
	SearchChatsTool = "search_chats"
)

const systemPrompt = `You are an assistant for college counselors. You answer questions about their students by searching the students' tasks and chat messages.

Choosing a tool:
- Use search_chats for anything a student or counselor said, discussed, mentioned, complained about, or felt.
- Use search_tasks for assignments, deadlines, due dates, progress, or work that is done, missing, or overdue.
- If the request is ambiguous, pick the tool that best fits the context of the request. Requests about opinions, feelings, or conversations go to search_chats. Requests about work, schedules, or progress go to search_tasks.
- Write the query as a short phrase describing what to find, keeping names and key topics from the request.

When you receive search results, summarize them in a few sentences. Only use what the results contain. If nothing relevant was found, say so plainly.`

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
