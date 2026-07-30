package shared

// HealthResponse is the standard health check response.
type HealthResponse struct {
	Status string `json:"status"`
}
