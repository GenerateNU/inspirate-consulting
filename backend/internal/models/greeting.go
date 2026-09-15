package models

type GreetingRequestBody struct {
	Name string `json:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}

type GreetingMessageBody struct {
	Message string `json:"greeting" example:"Hello, world!" doc:"Greeting message"`
}

type CreateGreetingInput struct {
	Body GreetingRequestBody
}

// This is the greet operation's output model, which has a body with a message.
type CreateGreetingOutput struct {
	Body GreetingMessageBody
}
