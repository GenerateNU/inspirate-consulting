package models

// Arguments for the search_tasks tool.
type SearchTasks struct {
	Query string `json:"query" doc:"What to search for in the student's tasks"`
}

// Arguments for the search_chats tool.
type SearchChats struct {
	Query string `json:"query" doc:"What to search for in the student's chats"`
}
