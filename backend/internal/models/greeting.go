package models

type MessageBody struct {
	Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
}

type CreateGreetingInput struct {
	// path tells huma to get name from the URL path
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

// This is the greet operation's output model, which has a body with a message.
type CreateGreetingOutput struct {
	Body MessageBody
}
