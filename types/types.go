package types

// Config holds the configuration for the checker.
type Config struct {
	Usernames struct {
		Custom bool `json:"custom"`
		Amount int  `json:"amount"`
		Length int  `json:"length"`
	} `json:"usernames"`

	Retry struct {
		Enabled     bool `json:"enabled"`
		MaxAttempts int  `json:"max_attempts"`
	} `json:"retry"`

	Threads int    `json:"threads"`
	Timeout int    `json:"timeout"`
	Webhook string `json:"webhook"`
}

// UsernameRequest is the request body for username checks.
type UsernameRequest struct {
	Username string `json:"username"`
}

// UsernameResponse is the response body for username checks.
type UsernameResponse struct {
	Taken bool `json:"taken"`
}
