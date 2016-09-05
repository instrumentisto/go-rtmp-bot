package model

// Value object of start test HTTP request.
type StartRequest struct {
	ServerURL   string `schema:"server"`       // RTMP media server URL.
	ModelCount  int    `schema:"model_count"`  // Count of model bots.
	ClientCount int    `schema:"client_count"` // Count of client bots.
}
